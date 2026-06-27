package server

import (
	"fmt"
	"log"
	"net"
)

type Server struct {
	cfg     Config
	ln      net.Listener
	handler HandlerFunc
}

type Config struct {
	Host    string
	Port    int
	Network string
}

type HandlerFunc func(net.Conn)

func (f HandlerFunc) handleClient(conn net.Conn) {
	f(conn)
}

func NewServer(cfg Config) *Server {
	return &Server{
		cfg: cfg,
	}
}

func (s *Server) HandleFunc(f HandlerFunc) {
	s.handler = f
}

func (s *Server) Run() error {
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	conn, err := net.Listen(s.cfg.Network, addr)
	if err != nil {
		return err
	}
	s.ln = conn
	return nil
}

func (s *Server) Close() error {
	return s.ln.Close()
}

func (s *Server) Start() {
	listener := s.ln
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}

		go s.handler.handleClient(conn)
	}
}
