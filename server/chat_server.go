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
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("[-] There should be one argument: Ex: 'localhost:9999'")
		return
	}

	fmt.Println("[*] Listening on " + args[1] + "...")

	listener, err := net.Listen("tcp", args[1])
	if err != nil {
		fmt.Println("[-] Error while listening on " + args[1] + ": " + err.Error())
		return
	}
	defer listener.Close()

	
	// Continously accept connections and handle data
	fmt.Println("[*] Accepting connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[-] Error accepting on" + args[1] + ": " + err.Error())
			fmt.Println("[*] Waiting for next connection...")
			return
		}

		go handleConnection(conn)
	}
}

// Handles the processsing of data from connection
// to be used as a goroutine
func handleConnection(conn net.Conn) {
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("[-] Error receiving data -> \t" + err.Error())
			if err.Error() == "EOF" {
				fmt.Println("\t: Ending goroutine...")
			} else {
				fmt.Println("\t: Unspecified communication error, ending goroutine...")
			}
			return
		}

		if string(data) == "halt\n" {
			fmt.Println("[*] Exiting...")
			os.Exit(1)
		}
		
		t := time.Now()
		timeReceived := t.Format(time.RFC3339)
		fmt.Print("[+] (" + timeReceived + ") " + string(data))
	}
}
