package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
	"net"
)

func ConnectToServer(ipaddr string, port string) {
	fmt.Println("[*] Starting client, connecting to tcp server [ip:port]: ", ipaddr, ":", port)

	ipWithPort := []string{ipaddr, port}
	
	connection, error := net.Dial("tcp", strings.Join(ipWithPort, ":"))
	if error != nil {
		fmt.Println("[-] Error connecting to server on port 8888")
	}
	// Reader processing data to write/read
	for {
		// Process input to send
		fmt.Print("Message to send: ")
		reader := bufio.NewReader(os.Stdin)
		text, error := reader.ReadString('\n')
		if error != nil {
			fmt.Println("[-] Error reading from stdin")
		}

		// Send to socket
		fmt.Fprintf(connection, text + "\n")
		
		// Listen for reply
		message, error := bufio.NewReader(connection).ReadString('\n')
		if error != nil {
			fmt.Println("[-] Error retrieving response from server")
		}
		fmt.Println("[+] Message from server: ", message)
	}
}
