package sysm

import (
	"encoding/json"
	"github.com/cloudfoundry/gosigar"
	"time"
)

const (
	WatcherInterval = 500
	DataChanSize    = 10
)

type Info struct {
	Cpu  float64 `json:"cpu"`
	Mem  float64 `json:"mem"`
	Time int64   `json:"time"`
}

type Watcher struct {
	ticker *time.Ticker
	stop   chan bool
	Data   chan *[]byte
}

func NewWatcher() *Watcher {
	return &Watcher{nil, nil, make(chan *[]byte)}
}

func (w *Watcher) Start() {
	go func() {
		w.ticker = time.NewTicker(time.Millisecond * WatcherInterval)
		defer w.ticker.Stop()
		w.stop = make(chan bool)
		for {
			select {
			case <-w.stop:
				return
			case <-w.ticker.C:
				var info Info
				cpu := sigar.Cpu{}
				cpu.Get()
				info.Cpu = float64(100) - float64(cpu.Idle*100)/float64(cpu.Total())

				mem := sigar.Mem{}
				mem.Get()
				info.Mem = float64(100) - float64(mem.Free)/float64(mem.Total)

				info.Time = time.Now().UnixNano() / 1000000

				data, _ := json.Marshal(info)
				w.Data <- &data
			}
		}
	}()
}

func (w *Watcher) Stop() {
	w.stop <- true
}
