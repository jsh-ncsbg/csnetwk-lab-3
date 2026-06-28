package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
)

type Server struct {
	cfg     Config
	ln      net.Listener
	ctx     context.Context
	handler HandlerFunc
}

type Config struct {
	Host    string
	Port    int
	Network string
}

type HandlerFunc func(context.CancelFunc, net.Conn)

func (f HandlerFunc) handleClient(ctxCancelFunc context.CancelFunc, conn net.Conn) {
	f(ctxCancelFunc, conn)
}

func NewServer(ctx context.Context, cfg Config) *Server {
	return &Server{
		ctx: ctx,
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

	if s.ctx == nil {
		s.ctx = context.Background()
	}

	done, cancel := context.WithCancel(s.ctx)
	defer cancel()

	go func() {
		<-done.Done()
		s.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				fmt.Println("Server listener closed. Stopping accept loop.")
				// TODO: replace with logger
				return
			}

			// TODO: replace with logger
			log.Print(err)
			continue
		}

		go s.handler.handleClient(cancel, conn)
	}
}
