package hub

import (
	"fmt"

	"github.com/papriwalprateek/engineer-chat/util"
)

// Settings stores the configuration settings for the server.
type Settings struct {
	Host string
	Port string
}

// Room stores the room details.
type Room struct {
	Clients []string
}

// Store is an in-memory storage for the running server.
var Store map[string]*Room

// ExecCommand executes the command.
func ExecCommand(action string, body string, client *Client) {
	var payload string

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
				Store[client.Room].RemoveClient(client)
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
					Store[client.Room].RemoveClient(client)
				}
				client.Room = body
				if _, ok := Store[client.Room]; !ok {
					Store[client.Room] = &Room{Clients: []string{client.Username}}
				} else {
					Store[client.Room].Clients = append(Store[client.Room].Clients, client.Username)
				}
				client.SendMessage("enter", body, false)
			} else {
				client.SendMessage("enter", "invlid room name", true)
			}

		// command to leave the given chat room
		case "leave":
			if client.Room != "lobby" {
				client.SendMessage("leave", client.Room, false)
				Store[client.Room].RemoveClient(client)
				client.Room = "lobby"
			} else {
				client.SendMessage("leave", "already in lobby", true)
			}

		// command to list the active rooms
		case "rooms":
			flag := false
			payload = "Active Rooms:"
			for r, room := range Store {
				if len(room.Clients) > 0 {
					flag = true
					payload += fmt.Sprintf("\n%v(%v)", r, len(room.Clients))
				}
			}
			if flag {
				client.SendMessage("rooms", payload, true)
			} else {
				client.SendMessage("rooms", "**no active room**", true)
			}

		// command to send private message to a user
		case "pm":
			rec, msg := util.ParseMsg("/" + body)
			payload = fmt.Sprintf("**pm** [%v] %v", client.Username, msg)
			client.SendPM(rec, payload)

		// command to list the users in the given room
		case "users":
			if body == "" {
				if client.Room == "lobby" {
					payload = "**users in lobby**"
				} else {
					payload = fmt.Sprintf("**users in [%v] room**", client.Room)
				}
				for _, cl := range ListClients() {
					if cl.Room == client.Room {
						payload += fmt.Sprintf("\n%v", cl.Username)
					}
				}
				client.SendMessage("users", payload, true)
			} else {
				if _, ok := Store[body]; ok {
					payload = fmt.Sprintf("**users in [%v] room**", body)
					for _, cl := range Store[body].Clients {
						payload += fmt.Sprintf("\n%v", cl)
					}
					client.SendMessage("users", payload, true)
				} else {
					client.SendMessage("users", "**no such room**", true)
				}
			}

		// command to list the available commands
		case "help":
			body = "**Engineer Chat**\n" +
				"Synopsis: /<command> <body>\n" +
				"List of Commands:\n" +
				"/register <username> : registers the user with the given username\n" +
				"/quit : logout\n" +
				"/message <message> : (or simply type your <message>) broadcast the message in the room\n" +
				"/pm <username> <message>: messages privately to the given user\n" +
				"/ignore <username> : ignores the user\n" +
				"/enter <room> : enters the given room\n" +
				"/leave : leave the room and come back in the lobby\n" +
				"/rooms : lists the available rooms\n" +
				"/users <room> : lists the users in the given room"
			client.SendMessage("help", body, true)

		default:
			client.SendMessage("unrecognized", action, true)
		}
	}
}

// RemoveClient removes client from the room.
func (room *Room) RemoveClient(client *Client) {
	clientsInRoom := room.Clients
	for i, c := range clientsInRoom {
		if c == client.Username {
			clientsInRoom = append(clientsInRoom[:i], clientsInRoom[i+1:]...)
			break
		}
	}
	room.Clients = clientsInRoom
}
