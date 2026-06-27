package client

import (
	"io"
	"net"
	"strconv"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
)

type Client struct {
	in        io.Reader
	conn      net.Conn
	serverCfg server.Config
	handler   HandlerFunc
}

type HandlerFunc func(io.Reader, io.ReadWriter)

func (f HandlerFunc) Serve(r io.Reader, rw io.ReadWriter) {
	f(r, rw)
}

func NewClient(in io.Reader, serverCfg server.Config) *Client {
	return &Client{
		in:        in,
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
	c.handler.Serve(c.in, c.conn)
}
