package main

import (
	"flag"
	"os"

	intlClient "github.com/jsh-ncsbg/csnetwk-lab-3/internal/client"
	pkgClient "github.com/jsh-ncsbg/csnetwk-lab-3/pkg/client"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/logger"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/server"
)

func main() {
	serverHost := flag.String("host", "localhost", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	flag.Parse()

	// logger setup
	logger := logger.NewLogger("CLIENT", os.Stdout)

	logger.Println("connecting to server...")
	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}

	client := pkgClient.NewClient(logger, os.Stdout, os.Stdin, serverConfig)
	client.HandleFunc(intlClient.Serve)

	if err := client.Dial(); err != nil {
		logger.Println("error connecting to server: ", err)
		os.Exit(1)
	}
	defer client.Close()
	logger.Println("connected to server")

	client.Start()
}
