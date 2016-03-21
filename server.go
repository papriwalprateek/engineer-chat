package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/papriwalprateek/engineer-chat/hub"
	"github.com/papriwalprateek/engineer-chat/util"
)

func main() {
	settings := &hub.Settings{Host: "localhost", Port: "5555"}
	ln, err := net.Listen("tcp", ":"+settings.Port)
	if err != nil {
		fmt.Println("Can't connect to server: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Chat server started on port %v...\n", settings.Port)

	// initialize hub Store
	hub.Store = make(map[string]*hub.Room)

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Can't accept connections: ", err.Error())
			os.Exit(1)
		}

		// keep track of the client details
		client := hub.Client{Connection: conn, Room: "lobby"}
		client.Register()

		// allow non-blocking client request handling
		channel := make(chan string)
		go waitForInput(channel, &client)
		go handleInput(channel, &client)

		client.SendMessage("login", "Welcome to the Engineer Chat Server!\nType /help to list the commands\nLogin Name?", true)
	}

}

// wait for client input (buffered by newlines) and signal the channel
func waitForInput(out chan string, client *hub.Client) {
	defer close(out)
	reader := bufio.NewReader(client.Connection)
	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			// connection has been closed, remove the client
			client.Close(true)
			return
		}
		out <- string(line)
	}
}

// listen for channel updates for a client and handle the message
func handleInput(in <-chan string, client *hub.Client) {
	for {
		message := <-in
		if message != "" {
			message = strings.TrimSpace(message)
			var action string
			var body string
			if message[0] == '/' {
				action, body = util.ParseMsg(message)
			} else {
				body = message
				if client.Username == "" {
					action = "register"
				} else {
					action = "message"
				}
			}
			hub.ExecCommand(action, body, client)
		}
	}
}
