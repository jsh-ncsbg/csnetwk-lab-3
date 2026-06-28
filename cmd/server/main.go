package main

import (
	"context"
	"flag"
	"os"

	intlServer "github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/server"
	pkgServer "github.com/jsh-ncsbg/csnetwk-lab-3/pkg/server"
)

func main() {
	serverHost := flag.String("host", "", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	// silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// logger setup
	logger := logger.NewLogger("SERVER", os.Stdin)

	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}

	server := pkgServer.NewServer(context.Background(), logger, serverConfig)

	err := server.Run()
	if err != nil {
		logger.Println("failed to start server", err)
		os.Exit(1)
	}
	defer server.Close()
	logger.Println("listening on ", serverConfig.Host, ":", serverConfig.Port)

	server.HandleFunc(intlServer.HandleClient)
	server.Start()
}
