package hub

import (
	"fmt"
	"net"
)

// Client stores the client details.
type Client struct {
	Connection net.Conn
	Username   string
	Room       string
	ignoring   []string
}

// in-memory storage
var clients []*Client

// Close the client connection and cleanup
func (client *Client) Close(doSendMessage bool) {
	if doSendMessage {
		// On sending the close command, the connection will terminate leading to another close
		// and hence this message will be sent.
		client.BroadcastMsg("disconnect", "")
	}
	client.Connection.Close()
	client.Delete()
}

// Register stores the client details.
func (client *Client) Register() {
	clients = append(clients, client)
}

// Ignore adds the user to the client's ignoring list.
func (client *Client) Ignore(user string) {
	client.ignoring = append(client.ignoring, user)
}

// IsIgnoring returns whether the client is ignoring the given user.
func (client *Client) IsIgnoring(username string) bool {
	for _, value := range client.ignoring {
		if value == username {
			return true
		}
	}
	return false
}

// BroadcastMsg broadcasts message to the clients.
func (client *Client) BroadcastMsg(msgType string, message string) {
	var payload string
	if msgType == "message" {
		payload = fmt.Sprintf("[%v] %v", client.Username, message)
	} else {
		payload = fmt.Sprintf("**%v** [%v] %v", msgType, client.Username, message)
	}

	for _, c := range clients {
		if (msgType == "message" && client.Room != c.Room) || c.IsIgnoring(client.Username) {
			continue
		}
		fmt.Fprintln(c.Connection, payload)
	}
}

// SendMessageToClientOnly sends message to this client only. This is mostly used to get
// information/help from the chat server.
func (client *Client) SendMessageToClientOnly(msgType string, message string) {
	var payload string
	if msgType == "unrecognized" {
		payload = fmt.Sprintf("**%v** unrecognized command", message)
	} else {
		payload = message
	}
	fmt.Fprintln(client.Connection, payload)
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

// Delete deletes the client entry from stored clients.
func (client *Client) Delete() {
	for i, c := range clients {
		if c.Username == client.Username {
			clients = append(clients[:i], clients[i+1:]...)
			break
		}
	}
}

// ListClients list all the clients in the hub.
func ListClients() []*Client {
	return clients
}
