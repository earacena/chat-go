// Project:      Serverside Go program for communication(chat/ftp)
// Author:       Emanuel Aracena
// Date created: August 19, 2019
// Name of file: chat_server.go
// Description : Handles server creation, incoming connections/messages, and
//               client activity.

package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
	"encoding/json"
	"io"
)

type Server struct {
	ipAndPort string
	listener  net.Listener
	conn      net.Conn
}

type Message struct {
	ID   string
	Body string
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There should be one argument: Ex: 'localhost:9999'")
		return
	}

	fmt.Println("[*] Listening on " + args[1] + "...")
	server := Server { ipAndPort: args[1], } 
	
	server.listen()
	defer server.listener.Close()

	// Continously accept connections and handle data
	fmt.Println("[*] Accepting connections...")
	
	for {
		server.acceptConnections()
		go server.HandleConnection(io.Writer(os.Stdin))
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
		fmt.Println("[-] Error while listening on " + s.ipAndPort + ": " + err.Error())
		return
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



