TCP Chat Server

A simple multi-client chat server written in Go. Clients connect over TCP, choose a username, and every message is broadcast to all other connected clients.

Features:

→ Multiple simultaneous clients

→ Usernames

→ Join and leave notifications

→ Messages are not echoed back to the sender

→ Empty messages are ignored

Requirements:
Go 1.xx or later

Running:
go run . [port]

The port defaults to 8080.

Connecting

From the same machine:

ncat localhost 8080

From another machine on the same network, use the server's IP address instead of localhost.

How it works
main listens for TCP connections and starts one goroutine per client.
handleConnection asks for a username, then reads the client's messages.
broadcast runs in a single goroutine, keeps the set of connected clients, and forwards messages to everyone except the sender.
clientWriter writes messages from a client's outgoing channel to its connection.