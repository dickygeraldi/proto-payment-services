package main

import (
	"fmt"
	"log"
	"os"
	cmd "proto-parking-services/pkg/cmd/server"
	"runtime"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

type Channel struct {
	Channel string `json:"channel"`
}

func sendJoin(c *gosocketio.Client) {
	log.Println("Acking /join")
	result, err := c.Ack("/join", Channel{"309241010"}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Ack result to /join: ", result)
	}
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())

	c, err := gosocketio.Dial(
		gosocketio.GetUrl(os.Getenv("SOCKET_HOST"), 80, true),
		transport.GetDefaultWebsocketTransport())

	if err != nil {
		log.Fatal("Error 1: ", err)
	}

	err = c.On("309241010", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})

	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal("Error 2 :", err)
	}

	time.Sleep(1 * time.Second)
	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")

	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
