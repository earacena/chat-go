package main

import (
	"bufio"
	"fmt"
	"strings"
	"net"
)

func CreateServerAndListen(ipaddr string, port string) {
	fmt.Println("[*] Starting tcp server on [ip address:port]: ", ipaddr, ":", port)

	// Join the ip address and port for a format readable by the Listen function
	
	ipWithPort := []string{ipaddr, port}
	listener, error := net.Listen("tcp", strings.Join(ipWithPort, ":"))
	if error != nil {
		fmt.Println("[-] Error creating tcp server on port 8888. ")
	}

	fmt.Println("[+] Server started successfully")
	
	defer listener.Close()
	// Continously wait and accept all incoming connections on port
	for {
		connection, error := listener.Accept()
		if error != nil {
			fmt.Println("[-] Error accepting connection on port 8888 [tcp]")
		}
		go handleConnection(connection)
	}
}

// Handles the data using a reader
func handleConnection(connection net.Conn) {
	message, error := bufio.NewReader(connection).ReadString('\n')
	if error != nil {
		fmt.Println("[-] Error receiving message.")
	}

	fmt.Print("[+] Data received: ", string(message))

	fmt.Fprintf(connection, "Your message was successfully received.")
}
