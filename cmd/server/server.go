package main

import (
	"log"
	"net/http"
	"os"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const SERVER_NAME = "SERVER"

func main() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")

		c.Emit("/message", Message{SERVER_NAME, "main", "using emit"})

		c.Join("test")
		c.BroadcastTo("test", "/message", Message{SERVER_NAME, "main", "using broadcast"})
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel Channel) string {
		time.Sleep(2 * time.Second)
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Printf("Serving at localhost:%s...\n", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), serveMux))
}
