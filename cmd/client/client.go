package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	
	socketio "github.com/graarh/golang-socketio"
)

const STOP_CHARACTER = "\r\n\r\n"

// func receive(conn net.Conn) <-chan string {

// }

func main() {
	addr := strings.Join([]string{os.Getenv("SERVER_IP"), os.Getenv("SERVER_PORT")}, ":")
	fmt.Printf("Connecting to %s...\n", addr)
	conn, err := net.Dial("tcp", addr)

	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}

	defer conn.Close()
	fmt.Printf("Connection: %v\n", conn)

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("-> ")

		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		fmt.Println(text)

		if strings.Compare("!quit", text) == 0 {
			log.Println("bye!")
			os.Exit(0)
		}

		conn.Write([]byte(text))
		conn.Write([]byte(STOP_CHARACTER))
	}
}
