package previewer

import (
	"fmt"
	"github.com/skratchdot/open-golang/open"
)

const (
	MarkdownChanSize = 3
	Version          = "0.1"
)

func NewPreviewer(port int) *Previewer {
	return &Previewer{port, nil, make(chan bool)}
}

type Previewer struct {
	port       int
	httpServer *HTTPServer
	stop       chan bool
}

func (p *Previewer) Run(files ...string) {
	p.httpServer = NewHTTPServer(p.port)
	p.httpServer.Listen()

	for _, file := range files {
		addr := fmt.Sprintf("http://localhost:%d/%s", p.port, file)
		open.Run(addr)
	}

	<-p.stop
}

func (p *Previewer) UseBasic() {
	MdConverter.UseBasic()
}
