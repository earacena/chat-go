// Project: Go Chat Application (Peer-to-Peer) communication
// Author: Emanuel Aracena Beriguete
// Date created: August 23, 2019
// Name of file: chat.go
// Description: Peer-to-Peer chat application using TCP.
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"encoding/json"
	"io"
	"strings"
)

type Server struct {
	ipAndPort string
	listener net.Listener
	conn net.Conn
}

type Client struct {
	ipAndPort string
	clientID string
	conn net.Conn
}

type Message struct {
	ID string
	Body string
}

func main() {

	fmt.Print("[?] Wait (server) or find (client) a connection " +
		"[server, client]: ")
	mode := ReadUserInput(bufio.NewReader(os.Stdin), false)
	mode = strings.Replace(mode, "\n", "", -1)

	var s Server
	var c Client

	if mode == "server" {
		fmt.Print("[?] Enter port for server: ")
		port := ReadUserInput(bufio.NewReader(os.Stdin), false)
		port = strings.Replace(port, "\n", "", -1)
		s = Server { ipAndPort: "localhost:" + port }
		s.listen()

		id := ChooseID(bufio.NewReader(os.Stdin), false)
		c = Client {
			clientID: id,
			ipAndPort: s.ipAndPort,	
		}

	}

	if mode == "client" {
		fmt.Print("[?] Enter IP and port to connect to: ")
		ip := ReadUserInput(bufio.NewReader(os.Stdin), false)
		id := ChooseID(bufio.NewReader(os.Stdin), false)

		ip = strings.Replace(ip, "\n", "", -1)
		id = strings.Replace(id, "\n", "", -1)
		
		fmt.Println("[*] Starting client, connecting to " + ip +
			" with ID " + id + "...")

		c = Client {
			clientID: id,
			ipAndPort: ip,
		}

		fmt.Print("[?] Enter a port for your own connection: ")
		port := ReadUserInput(bufio.NewReader(os.Stdin), false)
		port = strings.Replace(port, "\n", "", 1)
		s := Server { ipAndPort: "localhost:" + port }
		s.listen()

	}

	if mode != "server" && mode != "client" {
		fmt.Println("[-] Error: unrecognized response, try again... ")
		os.Exit(1)
	}
	
	defer s.listener.Close()	
	var err error
	peerFound := true
	for {
		err = c.ConnectToServer()
		if err != nil {
			fmt.Println("[-] No server to connect to, ensure peer" +
				" is ready...")
			peerFound = false
		}
		if err == nil && peerFound == false {
			fmt.Println("[+] Peer found.")
			peerFound = true
		}
		
		s.acceptConnections()
		go s.HandleConnection(io.Writer(os.Stdin))

		if peerFound == true {
			c.SendMessage(bufio.NewReader(os.Stdin))
		}
	}
}



// Handles the processsing of data from connection
// to be used as a goroutine
func (s *Server) HandleConnection(w io.Writer) {
	for {
		data := ReceiveMessage(bufio.NewReader(s.conn), false)
		m := DecodeMessage([]byte(data))
		frm, _ := FormatMessage(m)
		fmt.Fprint(w, frm)
	}
}

func (s *Server) listen() {
	var err error
	s.listener, err = net.Listen("tcp", s.ipAndPort)
	if err != nil {
		fmt.Println("[-] Error while listening on " + s.ipAndPort +
			": " + err.Error())
		os.Exit(1)
	}
}

func (s *Server) acceptConnections() {
	var err error
	s.conn, err = s.listener.Accept()
	if err != nil {
		fmt.Println("[-] Error accepting on" + s.ipAndPort + ": " + err.Error())
		fmt.Println("[*] Waiting for next connection...")
		return
	}
}

func ReceiveMessage(r *bufio.Reader, silent bool) string {
	data, err := r.ReadString('\n')
	if err != nil {
		fmt.Print("[-] Error receiving data -> \t" + err.Error() + "\n")
		if err.Error() == "EOF" {
			if silent != true {
				fmt.Println("\t: Client terminated connection, exiting...")
			}
			os.Exit(1)
		}

		if silent != true {
			fmt.Println("\t: Unspecified communication error, ending goroutine...")
		}
		return ""
	}
	
	return data
}

func DecodeMessage(encoded []byte) Message {
	var err error
	var m Message

	err = json.Unmarshal(encoded, &m)
	if err != nil {
		fmt.Println("[-] Error trying to unmarshal encoded message -> ", err.Error())
	}
	
	return m
}

func FormatMessage(m Message) (string, string) {
	t := time.Now()
	timeReceived := t.Format(time.RFC3339)
	return "[+] (" + timeReceived + ") " + m.ID + " " + m.Body + "\n", timeReceived
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
