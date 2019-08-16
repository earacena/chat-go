package client

import(
	"fmt"
	"net"
)

func ConnectToServer() {
     connection, error := net.Dial("tcp", ":8888")
     if err != nil {
     	fmt.Println("[-] Error connecting to server on port 8888")
     }
     // Reader processing data to write/read...
}