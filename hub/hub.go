package hub

import (
	"fmt"
	"net"
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

// Client stores the client details.
type Client struct {
	Connection net.Conn
	Username   string
	Room       string
	ignoring   []string
}

// in-memory storage
var clients []*Client

// Close the client connection and clenup
func (client *Client) Close(doSendMessage bool) {
	if doSendMessage {
		// if we send the close command, the connection will terminate causing another close
		// which will send the message
		client.SendMessage("disconnect", "", false)
	}
	client.Connection.Close()
	clients = removeEntry(client, clients)
}

// Register the connection and cache it
func (client *Client) Register() {
	clients = append(clients, client)
}

// Ignore adds the user to the client's ignoring list.
func (client *Client) Ignore(user string) {
	client.ignoring = append(client.ignoring, user)
}

// IsIgnoring returns whether the given user is in the client's ignoring list.
func (client *Client) IsIgnoring(username string) bool {
	for _, value := range client.ignoring {
		if value == username {
			return true
		}
	}
	return false
}

// SendMessage sends message to all clients
func (client *Client) SendMessage(messageType string, message string, thisClientOnly bool) {
	var payload string

	if thisClientOnly {
		if messageType == "unrecognized" {
			payload = fmt.Sprintf("**%v** unrecognized command", message)
		} else {
			payload = message
		}
		fmt.Fprintln(client.Connection, payload)
	} else if client.Username != "" {
		if messageType == "message" {
			payload = fmt.Sprintf("[%v] %v", client.Username, message)
		} else {
			payload = fmt.Sprintf("**%v** [%v] %v", messageType, client.Username, message)
		}

		for _, c := range clients {
			// write the message to the client
			if (thisClientOnly && c.Username == client.Username) || (!thisClientOnly && c.Username != "") {
				if messageType == "message" && client.Room != c.Room || c.IsIgnoring(client.Username) {
					continue
				}
				fmt.Fprintln(c.Connection, payload)
			}
		}
	}
}

// SendPM sends a private message to the receiver client.
func (client *Client) SendPM(rec string, msg string) {
	for _, cl := range clients {
		if cl.Username == rec {
			fmt.Fprintln(cl.Connection, msg)
			break
		}
	}
}

// Exists check whether the username is already occupied by some client.
func (client *Client) Exists(c string) bool {
	for _, cl := range clients {
		if c == cl.Username {
			return true
		}
	}
	return false
}

// remove client entry from stored clients
func removeEntry(client *Client, arr []*Client) []*Client {
	rtn := arr
	index := -1
	for i, value := range arr {
		if value == client {
			index = i
			break
		}
	}

	if index >= 0 {
		// we have a match, create a new array without the match
		rtn = make([]*Client, len(arr)-1)
		copy(rtn, arr[:index])
		copy(rtn[index:], arr[index+1:])
	}

	return rtn
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
