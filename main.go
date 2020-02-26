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
	MerchantApiKey string `json:"merchantApiKey"`
	Invoice        string `json:"invoice"`
	Status         string `json:"status"`
	Message        string `json:"message"`
	TrxId          string `json:"trxId"`
	Amount         string `json:"amount"`
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
		gosocketio.GetUrl(os.Getenv("SOCKET_HOST"), 80, false),
		transport.GetDefaultWebsocketTransport())

	if err != nil {
		log.Fatal("Error 1: ", err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal("Error 2 :", err)
	}

	go func() {
		err = c.On("309241010", func(h *gosocketio.Channel, args Message) {
			log.Println("--- Got chat message: ", args)
			log.Println("--- Got Data message: ", args.Invoice)
			log.Println("--- Got 1 message: ", args.MerchantApiKey)
			log.Println("--- Got 2 message: ", args.Message)
			log.Println("--- Got 3 message: ", args.TrxId)
			log.Println("--- Got 5 message: ", args.Status)
		})

		if err != nil {
			log.Fatal(err)
		}
	}()
	defer c.Close()

	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
