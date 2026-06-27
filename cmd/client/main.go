package main

import (
	"flag"
	"log"
	"os"

	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/client"
	"github.com/jsh-ncsbg/csnetwk-lab-3/internal/server"
)

func main() {
	serverHost := flag.String("host", "localhost", "server host name")
	serverPort := flag.Int("port", 8000, "server port number")

	flag.Parse()

	// logger setup
	log.SetPrefix("[CLIENT] ")
	log.SetFlags(log.Lmsgprefix | log.Ltime)

	log.Println("Connecting to server...")
	serverConfig := server.Config{
		Host:    *serverHost,
		Port:    *serverPort,
		Network: "tcp",
	}

	client := client.NewClient(os.Stdin, serverConfig)
	client.HandleFunc(serve)

	if err := client.Start(); err != nil {
		log.Fatalln("Error connecting to server: ", err)
	}

	log.Println("Connected to server")
	defer client.Close()
}
