package main

import (
	"fmt"
	"log"
	"net/http"

	"go-chat/pkg/models"
	"go-chat/pkg/utils"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const (
	SERVER_NAME = ""
	PORT        = "3000"
	CHANNEL     = "main"
	ROOM        = "main"
)

func disconectByChannel(users map[string]*gosocketio.Channel, c *gosocketio.Channel) {
	for name, channel := range users {
		if channel == c {
			delete(users, name)
		}
	}
}

func main() {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())
	var users map[string]*gosocketio.Channel = make(map[string]*gosocketio.Channel)

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("channel: ", c)
		log.Println("New User Connected")
	})
	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("channel: ", c)
		disconectByChannel(users, c)
		log.Println("Disconnected")
	})

	server.On("/join", func(c *gosocketio.Channel, channel models.Channel) models.Response {
		// time.Sleep(2 * time.Second)

		log.Println("channel: ", c)

		if _, ok := users[channel.From]; ok {
			return models.Response{true, fmt.Sprintf("User %s already exists. Try again", channel.From)}
		}

		msg := fmt.Sprintf("%s joined to %s\n", channel.From, channel.Channel)
		log.Println(msg)

		c.Join(ROOM)
		c.BroadcastTo(ROOM, "/notice", models.Message{utils.GetenvDefault("SERVER_NAME", SERVER_NAME), CHANNEL, msg})
		users[channel.From] = c

		return models.Response{false, "Joined to " + channel.Channel}
	})

	server.On("/message", func(c *gosocketio.Channel, message models.Message) models.Response {
		// time.Sleep(2 * time.Second)

		log.Println("channel: ", c)

		log.Println(message)
		c.BroadcastTo(ROOM, "/message", message)

		return models.Response{false, "Message broadcasted"}
	})

	server.On("/whisper", func(c *gosocketio.Channel, whisper models.Whisper) models.Response {
		// time.Sleep(2 * time.Second)

		log.Println("channel: ", c)

		log.Println(whisper)
		if _, ok := users[whisper.To]; ok {
			users[whisper.To].Emit("/whisper", whisper)
		} else {
			return models.Response{true, fmt.Sprintf("User %s not found. Search with !list command and try again", whisper.To)}
		}

		return models.Response{false, "Message whispered"}
	})

	server.On("/list", func(c *gosocketio.Channel) []string {
		// time.Sleep(2 * time.Second)

		log.Println("channel: ", c)

		list := make([]string, 0)
		for name, _ := range users {
			list = append(list, name)
		}
		return list
	})

	serveMux := http.NewServeMux()
	serveMux.Handle("/socket.io/", server)

	log.Printf("Serving at localhost:%s...\n", utils.GetenvDefault("PORT", PORT))
	log.Fatal(http.ListenAndServe(":"+utils.GetenvDefault("PORT", PORT), serveMux))
}
