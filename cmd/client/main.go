package main

import (
	"flag"
	"log"
	"os"

	intlClient "github.com/jsh-ncsbg/csnetwk-lab-3/internal/client"
	pkgClient "github.com/jsh-ncsbg/csnetwk-lab-3/pkg/client"
	"github.com/jsh-ncsbg/csnetwk-lab-3/pkg/server"
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

	client := pkgClient.NewClient(os.Stdin, serverConfig)
	client.HandleFunc(intlClient.Serve)

	if err := client.Dial(); err != nil {
		log.Fatalln("error connecting to server: ", err)
	}
	defer client.Close()
	log.Println("connected to server")

	client.Start()
}
