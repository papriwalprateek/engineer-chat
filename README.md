[![GoDoc](https://godoc.org/github.com/papriwalprateek/engineer-chat?status.png)](https://godoc.org/github.com/papriwalprateek/engineer-chat)

# engineer-chat

A simple chat server where clients(or rather engineers!) can easily chat through their terminal.

## Chat server

- Clone the repository in the `$GOPATH`
  ```
  $ go get github.com/papriwalprateek/engineer-chat
  ```

- Start the server
  ```
  $ go run server.go
  ```
  The server starts running at `localhost:5555`.

  Note: If you have [docker](https://www.docker.com/) installed on your machine, you can build the image and run the server inside the container.
  ```
  $ docker build -t engineer-chat .
  $ docker run --publish 6060:5555 --name test --rm engineer-chat
  ```
  This will start the server at `localhost:6060`.

  Alternatively, you can pull the latest image from docker hub and run it directly.
  ```
  $ docker run -d papriwalprateek/engineer-chat
  ```

## Chat clients
You can use either [telnet](http://linux.die.net/man/1/telnet) or [nc](http://linux.die.net/man/1/nc) to connect to the chat server.
```
$ telnet 127.0.0.1 5555
```

You can send commands or messages on this tcp stream. Commands begin with `/` and messages are anything else.

**List of Available Commands**

|Name           |Command               |Description
|---------------|----------------------|------------------------------------------- 
|Register       |`/register <username>`|registers the user with the given username.
|Logout         |`/quit`               |logout from the chat server.
|Message        |`/message <message>`  |Broadcast the message in the room. Also without `/`, automatically `/message` command is used.
|Private Message|`/pm <user> <message>`|messages privately to the given user.
|Join Room      |`/enter <room>`      |Joins the given room. If room is not there, it creates the room and then joins it.
|Leave Room     |`/leave`            |Leaves the current room and comes back into the lobby.
|List Rooms     |`/rooms`            |Lists the available rooms.
|List Users     |`/users <rooms`     |Lists the users in the given room. A simple `/users` will return the users present in the current room.
|Help           |`help`              |Lists the available commands.

#### Sample Chat Session
A sample chat session with 3 clients - `prateek`, `parijat`, `prateekp`

```
Trying ::1...
Connected to localhost.
Escape character is '^]'.
Welcome to the Engineer Chat Server!
Type /help to list the commands
Login Name?
prateek
**username already taken**
Login name?
prateekp
**connect** [prateekp] 
/users
**users in lobby**
prateek
parijat
prateekp
/pm parijat lets get into college room
**pm** [parijat] ok
**enter** [parijat] college
/users college
**users in [college] room**
parijat
/enter college
**enter** [prateekp] college
prateek cant here us because we are in **college** room
[prateekp] prateek cant here us because we are in **college** room
[parijat] yeah thats cool
/leave
**leave** [prateekp] college
/users college 
**users in [college] room**
parijat
/users
**users in lobby**
prateek
prateekp
**disconnect** [prateek] 
/users
**users in lobby**
prateekp
/quit
Connection closed by foreign host.
```
