package main

import (
	"io"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", "8080")
	if err != nil {
		log.Fatalf("Failed to bind to port 8080: %v", err)
	}

	defer listener.Close()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("Failed to accept incoming connection: %v", err)
			continue
		}
		go handleConnection(clientConn)
	}
}

func handleConnection ( clientConn net.Conn) {
	defer clientConn.Close()

	backendConn, err := net.Dial("tcp", "127.0.0.1:9090") 
	if err != nil {
		log.Printf("Failed to connect to backend target: %v", err)
		return
	}
	defer backendConn.Close()

	doneChan := make(chan struct{}, 2)

	go func() {
		_, _ = io.Copy(backendConn, clientConn)
		doneChan<- struct{}{}
	}()

	go func() {
		_, _ = io.Copy(clientConn, backendConn)
		doneChan<- struct{}{}
	}()

	<- doneChan
}
