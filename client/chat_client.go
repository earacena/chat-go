// Project:      Serverside Go program for communication(chat/ftp)
// Author:       Emanuel Aracena
// Date created: August 19, 2019
// Name of file: chat_client.go
// Description : Handles client connection to server and communication

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
	"io"
)

type Client struct {
	ipAndPort string
	clientID  string
	//writer     *bufio.Writer
	conn net.Conn
}

type Message struct {
	ID   string
	Body string
}

func main() {
	// Ensure program is given sufficient information
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There must be one args: Ex: localhost:9999")
		return
	}

	fmt.Println("[*] Starting client on " + args[1] + "...")
	ID := ChooseID(bufio.NewReader(os.Stdin), false)

	client := Client{
		clientID:  ID,
		ipAndPort: args[1],
	}

	var err error
	for {
		err = client.ConnectToServer()
		if err != nil {
			os.Exit(1)
		}
		
		client.SendMessage(bufio.NewReader(os.Stdin))
	}

}

// Connects to server
func (c *Client) ConnectToServer() error {
	var err error
	c.conn, err = net.Dial("tcp", c.ipAndPort)
	if err != nil {
		fmt.Println("[-] Error dialing " + c.ipAndPort + "->\n\t" + err.Error())
	}
	return err
}

// Handles reading and sending message
func (c *Client) SendMessage(r *bufio.Reader) {
	data := ReadUserInput(r, true)

	data = strings.Replace(string(data), "\n", "", -1)

	invoked := CheckHaltCommand(data, false)
	if invoked == true {
		os.Exit(1)
	}
	
	m := Message{
		ID:   c.clientID,
		Body: data,
	}

	b := EncodeMessage(m)	
	Send(c.conn, b)
}

func ReadUserInput(r *bufio.Reader, prompt bool) string {
	if prompt == true {
		fmt.Print("[Send]> ")
	}
	
	data, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("[-] Error reading string to send message -> " + err.Error())
		return ""
	}
	return data
}

// Check if Halt command used
func CheckHaltCommand(data string, silent bool) bool {
	if silent != true {
		if data == "!halt" {
			fmt.Println("[+] Quit command invoked, exiting...")
		}
	}
	return data == "!halt"
}

// Encode the Message struct in JSON
func EncodeMessage(m Message) []byte {
	b, err := json.Marshal(m)
	if err != nil {
		fmt.Println("[-] Error JSON encoding message -> " + err.Error())
		return nil
	}
	return b
}

// Send data through socket
func Send(w io.Writer, data []byte) error {
	// c.conn
	_, err := fmt.Fprintf(w, string(data)+"\n")
	if err != nil {
		fmt.Println("[-] Error sending message through socket -> " + err.Error())
		fmt.Println("[-] Connection with server terminated...")
	}
	return err
}

// Read user input and assign client ID
func ChooseID(r *bufio.Reader, silent bool) string {
	if silent != true {
		fmt.Print("[*] Please enter an ID (for example, user123): ")
	}
	ID, err := r.ReadString('\n')
	if err != nil {
		fmt.Println("[-] Error choosing an ID -> " + err.Error())
		fmt.Println("[+] Using default ID (default)...")
		ID = "default"
	}

	ID = "[" + strings.Replace(ID, "\n", "", -1) + "]"
	return ID
}
