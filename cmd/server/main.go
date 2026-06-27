package main

import (
	"flag"
	"io"
	"log"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
)

func main() {
	serverHost := flag.String("host", "", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	silent := flag.Bool("silent", false, "disable all log output")

	flag.Parse()

	// logger setup
	log.SetPrefix("[SERVER] ")
	log.SetFlags(log.Ltime | log.Lmsgprefix)
	if *silent {
		log.SetOutput(io.Discard)
	}

	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}
	server := server.NewServer(serverConfig)

	err := server.Run()
	if err != nil {
		log.Fatal("failed to start server", err)
	}
	defer server.Close()
	log.Print("listening on ", serverConfig.Host, ":", serverConfig.Port)

	server.HandleFunc(handleClient)
	server.Start()
}
