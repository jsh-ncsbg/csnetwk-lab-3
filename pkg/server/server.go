package server

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
)

type Server struct {
	cfg     Config
	ln      net.Listener
	ctx     context.Context
	logger  logger.Logger
	handler HandlerFunc
}

type Config struct {
	Host    string
	Port    int
	Network string
}

type HandlerFunc func(context.CancelFunc, logger.Logger, net.Conn)

func (f HandlerFunc) handleClient(ctxCancelFunc context.CancelFunc, logger logger.Logger, conn net.Conn) {
	f(ctxCancelFunc, logger, conn)
}

func NewServer(ctx context.Context, logger logger.Logger, cfg Config) *Server {
	return &Server{
		logger: logger,
		ctx:    ctx,
		cfg:    cfg,
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
				s.logger.Println("Server listener closed. Stopping accept loop.")
				return
			}

			s.logger.Println(err)
			continue
		}

		go s.handler.handleClient(cancel, s.logger, conn)
	}
}
