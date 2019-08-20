// Project:      Serverside Go program for communication(chat/ftp)
// Author:       Emanuel Aracena
// Date created: August 19, 2019
// Name of file: chat_client.go
// Description : Handles client connection to server and communication

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	ipAndPort  string
	clientID   string
	reader     *bufio.Reader
	//writer     *bufio.Writer
	conn       net.Conn
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There must be two arguments: Ex: localhost:9999")
		return
	}

	fmt.Print("Please enter an ID (for example, user123): ")
	ID, err := bufio.NewReader(os.Stdin).ReadString('\n')
	if err != nil {
		fmt.Println("[-] Error choosing an ID -> " + err.Error())
		fmt.Println("[+] Using default ID (user123)...")
		ID = "user123"
		return
	}

	// Remove newlines to prevent delimiting bugs
	ID = "[" + strings.Replace(ID, "\n", "", -1) + "]"
	
	client := Client{
		clientID:  ID,
		ipAndPort: args[1],
		reader:    bufio.NewReader(os.Stdin),
	}

	for {
		client.connectToServer(err)
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
	data, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Println("[-] Error reading string to send message -> " + err.Error())
		return
	}
	// Send data through socket
	data = c.clientID + " " + data + "\n" 
	fmt.Fprintf(c.conn, data)
}
