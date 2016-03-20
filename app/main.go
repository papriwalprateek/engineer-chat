package app

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/googollee/go-socket.io"
	"github.com/papriwalprateek/engineer-chat/util"
)

type client struct {
	name string
	room string
}

func init() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	var clients []*client

	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")
		so.Join(lobby)
		cl := &client{name: "Anon", room: lobby}
		so.On(chatChannel, func(msg string) {
			if msg != "" {
				msg = strings.TrimSpace(msg)
				var action string
				var body string
				var payload string
				if msg[0] == '/' {
					action, body = util.ParseMsg(msg)
				} else {
					action = "message"
				}

				if cl.name == "Anon" && action != "register" {
					so.Emit(chatChannel, "please register using /register <username>")
				} else {
					switch action {
					case cmdRegister:
						cl.name = body
						clients = append(clients, cl)
						payload = fmt.Sprintf("You are now registered as %v", cl.name)
						so.Emit(chatChannel, payload)
					case cmdMessage:
						so.Emit(chatChannel, msg)
						payload = fmt.Sprintf("[%v] %v", cl.name, msg)
						so.BroadcastTo(cl.room, chatChannel, payload)
					case cmdEnter:
						so.Leave(cl.room)
						so.Join(body)
						cl.room = body
						payload = fmt.Sprintf("**%v** joins **%v** room", cl.name, cl.room)
						so.Emit(chatChannel, payload)
						so.BroadcastTo(lobby, chatChannel, payload)
					case cmdLeave:
						so.Leave(cl.room)
						so.Join(lobby)
						payload = fmt.Sprintf("**%v** leaves **%v** room", cl.name, cl.room)
						so.Emit(chatChannel, payload)
						so.BroadcastTo(lobby, chatChannel, payload)
						cl.room = lobby
					default:
						payload = fmt.Sprintf("**%v** command not supported", action)
						so.Emit(chatChannel, payload)
					}
				}
			}
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	http.Handle("/socket.io/", server)
	http.Handle("/", http.FileServer(http.Dir("./asset")))
	http.ListenAndServe(":8080", nil)
}
