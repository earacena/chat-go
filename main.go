package main

import (
	"fmt"
)

func main() {
	fmt.Println("Running Server and Client routines...")
	go CreateServerAndListen("localhost", "8888")
	go ConnectToServer("localhost", "8888")
}
