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

	listener, err := net.Listen("tcp", args[1])
	if err != nil {
		fmt.Println("[-] Error while listening on " + args[1] + ": " + err.Error())
		return
	}
	defer listener.Close()

	
	// Continously accept connections and handle data
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("[-] Error while accepting connection on " + args[1] + ": " +
				err.Error())
			return
		}
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("[-] Error while listenting on " + args[1] + err.Error())
			return
		}

		t := time.Now()
		timeReceived := t.Format(time.RFC3339)
		fmt.Print("[+] (" + timeReceived + "): " + string(data))
	}
}
