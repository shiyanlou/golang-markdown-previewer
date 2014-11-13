package previewer

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	ListeningTestInterval = 500
	MaxListeningTestCount = 10
)

type HTTPServer struct {
	port     int
	listener net.Listener
}

func NewHTTPServer(port int) *HTTPServer {
	return &HTTPServer{port, nil}
}

func (s *HTTPServer) Addr() string {
	return ":" + strconv.Itoa(s.port)
}

func (s *HTTPServer) ListenAndServe() {
	var err error
	server := &http.Server{
		Addr:           s.Addr(),
		Handler:        s,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.listener, err = net.Listen("tcp", s.Addr())
	if err != nil {
		panic(err)
	}

	server.Serve(s.listener)
}

func (s *HTTPServer) Listen() {
	go s.ListenAndServe()

	isListening := make(chan bool)
	go func() {
		result := false
		ticker := time.NewTicker(time.Millisecond * ListeningTestInterval)
		for i := 0; i < MaxListeningTestCount; i++ {
			<-ticker.C
			resp, err := http.Get("http://localhost" + s.Addr() + "/ping")
			if err == nil && resp.StatusCode == 200 {
				result = true
				break
			}
		}
		ticker.Stop()
		isListening <- result
	}()

	if <-isListening {
		fmt.Println("Listening", s.Addr(), "...")
	} else {
		panic("Can't connect to server")
	}
}

func (s *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path[1:] // remove '/'
	if path == "ping" {
		w.Write([]byte("pong"))
		fmt.Println("accept conn ...")
	} else if isWebsocketRequest(r) {
		NewWebsocket(path).Serve(w, r)
	} else {
		if strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown") {
			Template(w, path, s.port)
		} else {
			s.ServeStatic(w, path)
		}
	}
}

func (s *HTTPServer) ServeStatic(w http.ResponseWriter, path string) {
	if stat, err := os.Stat(path); err == nil && stat.Mode().IsRegular() {
		file, _ := os.Open(path)
		io.Copy(w, file)
	}
}

func contains(arr []string, needle string) bool {
	for _, v := range arr {
		if strings.Contains(v, needle) {
			return true
		}
	}
	return false
}

func isWebsocketRequest(r *http.Request) bool {
	upgrade := r.Header["Upgrade"]
	connection := r.Header["Connection"]
	return contains(upgrade, "websocket") && contains(connection, "Upgrade")
}

func (s *HTTPServer) Stop() {
	s.listener.Close()
}
