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
				action, body = util.ParseMsg(message)
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
					if client.Exists(body) {
						client.SendMessage("warn", "**username already taken**\nLogin name?", true)
					} else {
						client.Username = body
						client.SendMessage("connect", "", false)
					}

				// command to logout of the chat server
				case "quit":
					if client.Room != "lobby" {
						hubStore[client.Room].RemoveClient(client)
					}
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
					rec, msg := util.ParseMsg("/" + body)
					fmt.Println(rec, msg)
					payload := fmt.Sprintf("**pm** [%v] %v", client.Username, msg)
					client.SendPM(rec, payload)

				case "help":
					body = "**Engineer Chat**\n" +
						"Synopsis: /<command> <body>\n" +
						"List of Commands:\n" +
						"/register <username> : registers with the given username\n" +
						"/message <message> : (or simply type your <message>) broadcast the message in the room\n" +
						"/quit : logout\n" +
						"/ignore <username> : ignores the user\n" +
						"/enter <room> : enters the given room\n" +
						"/leave : leave the room and come back in the lobby\n" +
						"/rooms : lists the available rooms\n" +
						"/pm <username> <message>: messages privately to the given user"
					client.SendMessage("help", body, true)

				default:
					client.SendMessage("unrecognized", action, true)
				}
			}
		}
	}
}
