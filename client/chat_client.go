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
	"encoding/json"
)

type Client struct {
	ipAndPort  string
	clientID   string
	reader     *bufio.Reader
	//writer     *bufio.Writer
	conn       net.Conn
}

type Message struct {
	ID   string
	Body string
}

func main() {
	// Ensure program is given sufficient information
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There must be two arguments: Ex: localhost:9999")
		return
	}

	fmt.Println("[*] Starting client on " + args[1] + "...")
	
	// Initialize ID from user input
	fmt.Print("[*] Please enter an ID (for example, user123): ")
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

	// Continously perform connection to server and send data as needed
	for {
		client.connectToServer(err)
		client.sendMessage()
	}

}

// Connects to server
func (c *Client) connectToServer(err error) {
	c.conn, err = net.Dial("tcp", c.ipAndPort)
	if err != nil {
		fmt.Println("[-] Error dialing to " + c.ipAndPort + ": " + err.Error())
	}
}

// Handles reading and sending message
func (c *Client) sendMessage() {
	// Read user input
	fmt.Print("[Send]> ")
	data, err := c.reader.ReadString('\n')
	if err != nil {
		fmt.Println("[-] Error reading string to send message -> " + err.Error())
		return
	}

	// Close client using command "!halt"
	data = strings.Replace(string(data), "\n", "", -1)
	if data == "!halt" {
		fmt.Println("[+] Quit command invoked, exiting...")
		os.Exit(1)
	}

	// Create Message used in JSON
	m := Message{
		ID:   c.clientID,
		Body: data,
	}

	// Convert to JSON format
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("[-] Error JSON encoding message -> " + err.Error())
		return
	}
	// Send data through socket
	// Convert JSON []byte into string for simplicity, then server converts back to []byte
	_, err = fmt.Fprintf(c.conn, string(b) + "\n")
	if err != nil {
		fmt.Println("[-] Error sending message through socket -> " + err.Error())
		return
	}

	//fmt.Println("[+] Number of bytes send: ", bytesWritten)
}

