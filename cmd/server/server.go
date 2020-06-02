package main

import (
	"log"
	"net/http"
	"time"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"github.com/marcosrachid/go-chat/internal/models"
	"github.com/marcosrachid/go-chat/internal/utils"
)

const (
	SERVER_NAME = "SERVER"
	PORT        = "3000"
)

func main() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("Connected")

		c.Emit("/message", models.Message{utils.GetenvDefault("SERVER_NAME", SERVER_NAME), "main", "using emit"})

		c.Join("test")
		c.BroadcastTo("test", "/message", models.Message{utils.GetenvDefault("SERVER_NAME", SERVER_NAME), "main", "using broadcast"})
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel models.Channel) string {
		time.Sleep(2 * time.Second)
		log.Println("Client joined to ", channel.Channel)
		return "joined to " + channel.Channel
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Printf("Serving at localhost:%s...\n", utils.GetenvDefault("PORT", PORT))
	log.Fatal(http.ListenAndServe(":"+utils.GetenvDefault("PORT", PORT), serveMux))
}
