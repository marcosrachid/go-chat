package main

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go-chat/pkg/models"
	"go-chat/pkg/utils"

	gosocketio "github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
)

const (
	SERVER_IP   = "localhost"
	SERVER_PORT = "3000"
	CHANNEL     = "main"
)

func showHelp() {
	fmt.Println(
		`
Command usage:

	<command> [arguments]

The commands are:

	!list		list connected users
	!whisper	whisper to a specific user
	!quit		to exit
	!help		shows help
		`,
	)
}

func sendJoin(name string, c *gosocketio.Client) bool {
	log.Println("Joining channel...")
	r, err := c.Ack("/join", models.Channel{name, CHANNEL}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	var response models.Response
	json.Unmarshal([]byte(r), &response)
	if response.Error {
		fmt.Print(response.Message + ": ")
	}
	return response.Error
}

func sendMessage(name, message string, c *gosocketio.Client) {
	_, err := c.Ack("/message", models.Message{name, CHANNEL, message}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
}

func sendWhisper(name, to, message string, c *gosocketio.Client) {
	r, err := c.Ack("/whisper", models.Whisper{models.Message{name, CHANNEL, message}, to}, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	var response models.Response
	json.Unmarshal([]byte(r), &response)
	if response.Error {
		fmt.Println("\n" + response.Message)
		fmt.Print(">> ")
	}
}

func sendList(c *gosocketio.Client) {
	r, err := c.Ack("/list", nil, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	var response []string
	json.Unmarshal([]byte(r), &response)
	fmt.Println("User list: \n")
	for _, value := range response {
		fmt.Println(value)
	}
	fmt.Println("")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	ip := utils.GetenvDefault("SERVER_IP", SERVER_IP)
	port, _ := strconv.Atoi(utils.GetenvDefault("SERVER_PORT", SERVER_PORT))
	c, err := gosocketio.Dial(
		gosocketio.GetUrl(ip, port, false),
		transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	err = c.On("/notice", func(h *gosocketio.Channel, message models.Message) {
		log.Printf(" %s", message.Text)
		fmt.Print(">> ")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *gosocketio.Channel, message models.Message) {
		log.Printf(" %s: %s", message.From, message.Text)
		fmt.Print(">> ")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/whisper", func(h *gosocketio.Channel, whisper models.Whisper) {
		log.Printf(" %s whispered to you: %s", whisper.From, whisper.Text)
		fmt.Print(">> ")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Bye!")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	fmt.Print("What's your name? ")
	var name string
	for {
		time.Sleep(1 * time.Second)
		name = utils.ReadInput()
		if !sendJoin(name, c) {
			break
		}
	}

	showHelp()

Loop:
	for {
		fmt.Print(">> ")
		text := utils.ReadInput()
		switch {
		case strings.Compare(text, "!list") == 0:
			sendList(c)
		case strings.HasPrefix(text, "!whisper"):
			r := regexp.MustCompile("[^\\s]+")
			splitedText := r.FindAllString(text, -1)
			msg := strings.Join(splitedText[2:], " ")
			if len(msg) > 0 {
				go sendWhisper(name, splitedText[1], msg, c)
			}
		case strings.Compare(text, "!quit") == 0:
			break Loop
		case strings.Compare(text, "!help") == 0:
			showHelp()
		case strings.HasPrefix(text, "!"):
			r := regexp.MustCompile("[^\\s]+")
			splitedText := r.FindAllString(text, -1)
			fmt.Printf("Command %s not found\n", splitedText[0])
		case len(text) > 0:
			go sendMessage(name, text, c)
		default:
		}
	}
}
