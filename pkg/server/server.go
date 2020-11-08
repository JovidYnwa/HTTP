package server

import (
	"bytes"
	"io"
	"log"
	"net"
	"strings"
	"sync"
)

type HadlelerFunc func(conn net.Conn)

type Server struct {
	addr     string
	mu       sync.RWMutex
	handlers map[string]HadlelerFunc
}

func NewSever(addr string) *Server {
	return &Server{addr: addr, handlers: make(map[string]HadlelerFunc)}
}

func (s *Server) Register(path string, handler HadlelerFunc) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers[path] = handler
}

func (s *Server) Start() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Print(err)
		return err
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handle(conn)

	}
}

func (s *Server) handle(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {

			log.Print(err)
		}
	}()

	buf := make([]byte, 4096)
	n, err := conn.Read(buf)
	if err == io.EOF {
		log.Printf("%s", buf[:n])
	}

	data := buf[:n]

	requestLineDelim := []byte{'\r', '\n'}
	requestLineEnd := bytes.Index(data, requestLineDelim)
	if requestLineEnd == -1 {
		return
	}

	requestLine := string(data[:requestLineEnd])
	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		log.Print("length: ", parts)
		return
	}
	s.mu.RLock()

	if hadler, ok := s.handlers[parts[1]]; ok {
		s.mu.RUnlock()

		hadler(conn)
	}
	return

}

func Some() string {
	return "hey"
}
