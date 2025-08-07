package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// Server
func main() {
	listener, err := net.Listen("tcp", ":9991")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Server listening on :9991")

	conn, err := listener.Accept()
	if err != nil {
		return
	}
	fmt.Println("Proxy connected")

	go func(c net.Conn) {
		defer c.Close()
		reader := bufio.NewReader(c)
		for {
			msg, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			log.Printf("Received message: %s\n", msg)
			msg = strings.TrimSpace(msg)
			response := fmt.Sprintf("Hello %s", msg)
			c.Write([]byte(response + "\n"))
		}
	}(conn)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan
}
