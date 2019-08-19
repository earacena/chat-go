// Project:      Serverside Go program for communication(chat/ftp)
// Author:       Emanuel Aracena
// Date created: August 19, 2019
// Name of file: chat_client.go
// Description : Handles client connection to server and communication

package main

import (
	"fmt"
	"bufio"
	"os"
	"net"
)

type Client struct {
	ipAndPort string
	reader     *bufio.Reader
	//writer     *bufio.Writer
	conn       net.Conn
	messageIn  chan *string
	messageOut chan *string
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There must be two arguments: Ex: localhost:9999")
		return
	} 

	client := Client {
		ipAndPort: args[1],
		reader:    bufio.NewReader(os.Stdin),
	}

	var err error
	for {
		client.connectToServer(err)	
		// Continously send messages to server
		client.sendMessage()
	}

}

func (c *Client) connectToServer(err error) {
	c.conn, err = net.Dial("tcp", c.ipAndPort)
	if err != nil {
		fmt.Println("[-] Error dialing to " + c.ipAndPort + ": " + err.Error())
	}
}

// Handles reading and sending message
func (c *Client) sendMessage() {
	fmt.Print("[Send]> ")
	data, _ := c.reader.ReadString('\n')
	// Send data through socket
	fmt.Fprintf(c.conn, data + "\n")
}

// Handles receiving messages
func (c *Client) receiveMessage() {

}
