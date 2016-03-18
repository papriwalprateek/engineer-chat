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

// Client stores the client details.
type Client struct {
	// the client's connection
	Connection net.Conn
	// the client's username
	Username string
	// the current room or "global"
	Room string
	// list of usernames we are ignoring
	ignoring []string
	// the config properties
	Settings *Settings
}

// static client list
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

	if thisClientOnly {
		// this message is only for the provided client
		message = fmt.Sprintf("/%v", messageType)
		fmt.Fprintln(client.Connection, message)

	} else if client.Username != "" {
		// construct the payload to be sent to clients
		payload := fmt.Sprintf("/%v [%v] %v", messageType, client.Username, message)

		for _, _client := range clients {
			// write the message to the client
			if (thisClientOnly && _client.Username == client.Username) ||
				(!thisClientOnly && _client.Username != "") {

				// you should only see a message if you are in the same room
				if messageType == "message" && client.Room != _client.Room || _client.IsIgnoring(client.Username) {
					continue
				}

				// you won't hear any activity if you are anonymous unless thisClientOnly
				// when current client will *only* be messaged
				fmt.Fprintln(_client.Connection, payload)
			}
		}
	}
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
