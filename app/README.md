# webapp

The webapp is running at https://engineer-chat.appspot.com/. The app uses [socket.io](http://socket.io/) to maintain persistent connections.

Note: This app supports a limited set of commands.
- `/register <username>`
- `/message <message>` 
- `/enter <room>`
- `/leave`

## TODO
- Think of an architecture to allow communication between telnet clients and webapp clients. tcp <-> socket communication.
- Support more commands.
- Improve UI.
