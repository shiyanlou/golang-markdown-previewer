package sysm

type Sysm struct {
	port       int
	httpServer *HTTPServer
	stop       chan bool
}

func NewSysm(port int) *Sysm {
	return &Sysm{port, nil, make(chan bool)}
}

func (s *Sysm) Run() {
	s.httpServer = NewHTTPServer(s.port)
	s.httpServer.Listen()
	<-s.stop
}

func (s *Sysm) Stop() {
	s.httpServer.Stop()
	s.stop <- true
}
