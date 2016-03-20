package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"regexp"
	"strings"

	"github.com/papriwalprateek/engineer-chat/hub"
)

// hubStore is an in-memory storage for the running server.
var hubStore map[string]*hub.Room

func main() {
	settings := &hub.Settings{Host: "localhost", Port: "5555"}
	ln, err := net.Listen("tcp", ":"+settings.Port)
	if err != nil {
		fmt.Println("Can't connect to server: ", err.Error())
		os.Exit(1)
	}
	fmt.Printf("Chat server started on port %v...\n", settings.Port)

	// initialize hubStore
	hubStore = make(map[string]*hub.Room)

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

		client.SendMessage("login", "Welcome to the Engineer Chat Server! \nLogin Name?", true)
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
// messages must be in the format of /{action} {content} where content is optional depending on the action
// supported actions are "user", "chat", and "quit".  the "user" must be set before any chat messages are allowed
func handleInput(in <-chan string, client *hub.Client) {

	for {
		message := <-in
		if message != "" {
			message = strings.TrimSpace(message)
			var action string
			var body string
			if message[0] == '/' {
				action, body = parse(message)
			} else {
				body = message
				if client.Username == "" {
					action = "register"
				} else {
					action = "message"
				}
			}

			if action != "" {
				switch action {

				// command to post message on chat server
				case "message":
					client.SendMessage("message", body, false)

					// command to login into server by providing username
				case "register":
					client.Username = body
					client.SendMessage("connect", "", false)

				// command to logout of the chat server
				case "quit":
					hubStore[client.Room].RemoveClient(client)
					client.Close(false)

				// command to add the given username to client's ignoring list
				case "ignore":
					client.Ignore(body)
					client.SendMessage("ignoring", body, false)

				// command to enter the given chat room
				case "enter":
					if body != "" {
						if client.Room != "lobby" {
							hubStore[client.Room].RemoveClient(client)
						}
						client.Room = body
						if _, ok := hubStore[client.Room]; !ok {
							hubStore[client.Room] = &hub.Room{Clients: []string{client.Username}}
						} else {
							hubStore[client.Room].Clients = append(hubStore[client.Room].Clients, client.Username)
						}
						client.SendMessage("enter", body, false)
					}

				// command to leave the given chat room
				case "leave":
					if client.Room != "lobby" {
						client.SendMessage("leave", client.Room, false)
						hubStore[client.Room].RemoveClient(client)
						client.Room = "lobby"
					}

				case "rooms":
					body = "Available Rooms:"
					for r, room := range hubStore {
						body += fmt.Sprintf("\n%v(%v)", r, len(room.Clients))
					}
					client.SendMessage("rooms", body, true)

				case "pm":
					rec, msg := parse("/" + body)
					fmt.Println(rec, msg)
					payload := fmt.Sprintf("**pm** [%v] %v", client.Username, msg)
					client.SendPM(rec, payload)
				default:
					client.SendMessage("unrecognized", action, true)
				}
			}
		}
	}
}

// parse out message contents (/{action} {message})
func parse(message string) (string, string) {
	actionRegex, _ := regexp.Compile(`^\/([^\s]*)\s*(.*)$`)
	res := actionRegex.FindAllStringSubmatch(message, -1)
	if len(res) == 1 {
		return res[0][1], res[0][2]
	}
	return "", ""
}
