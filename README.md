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
  The server starting running at `localhost:5555`.

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

## Chat Clients
You can use either [telnet](http://linux.die.net/man/1/telnet) or [nc](http://linux.die.net/man/1/nc) to connect to the chat server.
```
$ telnet 127.0.0.1 5555
```

You can send commands or messages on this tcp stream. Commands begin with `/` and messages are anything else.

**List of Available Commands**
/register <username> : registers with the given username
/quit : logout
/message <message> : (or simply type your <message>) broadcast the message in the room
/pm <username> <message>: messages privately to the given user/ignore <username> : ignores the user
/enter <room> : enters the given room
/leave : leave the room and come back in the lobby
/rooms : lists the available rooms
/users <room> : lists the users in the given room
/help : list the available commands
