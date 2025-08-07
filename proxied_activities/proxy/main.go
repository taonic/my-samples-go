package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"sync"
)

const (
	proxyPort  = ":9992"
	serverAddr = "localhost:9991"
)

type Proxy struct {
	serverConn net.Conn
	mu         sync.Mutex
	respChan   chan string
}

func main() {
	proxy := &Proxy{
		respChan: make(chan string, 1),
	}

	if err := proxy.connectToServer(); err != nil {
		log.Fatal("Failed to connect to server:", err)
	}
	defer proxy.serverConn.Close()

	//go proxy.handleServerResponses()

	listener, err := net.Listen("tcp", proxyPort)
	if err != nil {
		log.Fatal("Failed to start proxy:", err)
	}
	defer listener.Close()

	fmt.Printf("Proxy listening on %s, forwarding to %s\n", proxyPort, serverAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go proxy.handleClient(clientConn)
	}
}

func (p *Proxy) connectToServer() error {
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		return err
	}
	p.serverConn = conn
	return nil
}

func (p *Proxy) handleClient(clientConn net.Conn) {
	defer func() {
		fmt.Println("Worker disconnected")
		clientConn.Close()
	}()

	fmt.Println("Worker connected")

	reader := bufio.NewReader(clientConn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		p.mu.Lock()
		p.serverConn.Write([]byte(msg))
		resp, err := bufio.NewReader(p.serverConn).ReadString('\n')
		if err != nil {
			return
		}
		p.mu.Unlock()

		clientConn.Write([]byte(resp))
	}
}
