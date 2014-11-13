package previewer

import (
	"io/ioutil"
	"os"
	"time"
)

const (
	WatcherInterval = 500
	DataChanSize    = 10
)

type DataChan struct {
	Raw chan *[]byte
	Req chan bool
}

type Watcher struct {
	path   string
	ticker *time.Ticker
	stop   chan bool
	C      *DataChan
}

func NewWatcher(path string) *Watcher {
	dataChan := DataChan{make(chan *[]byte, DataChanSize), make(chan bool)}
	return &Watcher{path, nil, nil, &dataChan}
}

func (w *Watcher) Start() {
	go func() {
		w.ticker = time.NewTicker(time.Millisecond * WatcherInterval)
		defer w.ticker.Stop()
		w.stop = make(chan bool)
		var currentTimestamp int64
		for {
			select {
			case <-w.stop:
				return
			case <-w.ticker.C:
				reload := false
				select {
				case <-w.C.Req:
					reload = true
				default:
				}

				info, err := os.Stat(w.path)
				if err != nil {
					continue
				}

				timestamp := info.ModTime().Unix()
				if currentTimestamp < timestamp || reload {
					currentTimestamp = timestamp

					raw, err := ioutil.ReadFile(w.path)
					if err != nil {
						continue
					}

					w.C.Raw <- &raw
				}
			}
		}
	}()
}

func (w *Watcher) Stop() {
	w.stop <- true
}
