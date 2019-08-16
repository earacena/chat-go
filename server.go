package server

import (
       "fmt"
       "net"
)

func CreateServerAndListen() {
     listener, error := net.Listen("tcp", ":8888")
     if error != nil {
     	fmt.Println("[-] Error creating tcp server on port 8888. ")
     }
     for {
     	 connection, error := listener.Accept()
	 if error != nil {
	    fmt.Println("[-] Error accepting connection on port 8888 [tcp]")
	 }
	 go handleConnection(connection)
     }
}

func handleConnection(conn Conn) {
     // protocol for data transmitted
}