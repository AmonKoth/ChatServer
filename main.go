package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type client struct {
	name     string
	outgoing chan string
}

type message struct {
	sender *client
	text   string
}

var (
	entering = make(chan *client)
	leaving  = make(chan *client)
	messages = make(chan message)
)

func main() {

	port := "8080"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		WriteupError(err)
		os.Exit(1)
	}
	defer listener.Close()
	go broadcast()

	for {
		connection, err := listener.Accept()
		if err != nil {
			WriteupError(err)
			continue
		}
		go handleConnection(connection)

	}

}

func handleConnection(connection net.Conn) {
	defer connection.Close()

	currentClient := &client{
		name:     "default name",
		outgoing: make(chan string, 10),
	}

	scanner := bufio.NewScanner(connection)
	var name string
	for name == "" {
		fmt.Fprint(connection, "Enter your name: ")
		if !scanner.Scan() {
			return
		}
		name = strings.TrimSpace(scanner.Text())
	}
	currentClient.name = name

	go clientWriter(connection, currentClient.outgoing)

	messages <- message{sender: currentClient, text: "- " + currentClient.name + " joined the channel"}
	entering <- currentClient

	for scanner.Scan() {
		text := strings.TrimSpace(scanner.Text())
		if text == "" {
			continue
		}
		messages <- message{sender: currentClient, text: currentClient.name + ": " + text}
	}

	messages <- message{sender: currentClient, text: "- " + currentClient.name + " left the channel"}
	leaving <- currentClient
}

func broadcast() {
	clients := make(map[*client]bool)

	for {
		select {
		case currentClient := <-entering:
			clients[currentClient] = true
		case currentClient := <-leaving:
			delete(clients, currentClient)
			close(currentClient.outgoing)
		case incoming := <-messages:
			for currentClient := range clients {
				if currentClient != incoming.sender {
					currentClient.outgoing <- incoming.text
				}
			}
		}
	}
}

func clientWriter(connection net.Conn, channel <-chan string) {
	for message := range channel {
		fmt.Fprint(connection, message+"\r\n")
	}
}
