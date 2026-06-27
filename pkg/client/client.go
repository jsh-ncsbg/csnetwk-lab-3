package client

import (
	"io"
	"net"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/server"
)

type Client struct {
	in        io.Reader
	out       io.Writer
	logger    logger.Logger
	conn      net.Conn
	serverCfg server.Config
	handler   HandlerFunc
}

type HandlerFunc func(logger.Logger, io.Writer, io.Reader, io.ReadWriter)

func (f HandlerFunc) Serve(logger logger.Logger, w io.Writer, r io.Reader, conn io.ReadWriter) {
	f(logger, w, r, conn)
}

func NewClient(logger logger.Logger, out io.Writer, in io.Reader, serverCfg server.Config) *Client {
	return &Client{
		in:        in,
		out:       out,
		logger:    logger,
		serverCfg: serverCfg,
	}
}

func (c *Client) HandleFunc(f HandlerFunc) {
	c.handler = f
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Dial() error {
	addr := net.JoinHostPort(c.serverCfg.Host, strconv.Itoa(c.serverCfg.Port))

	conn, err := net.Dial(c.serverCfg.Network, addr)
	if err != nil {
		return err
	}

	c.conn = conn
	return nil
}

func (c *Client) Start() {
	c.handler.Serve(c.logger, c.out, c.in, c.conn)
}
