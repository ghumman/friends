package main

import (
	"io"
	"net"
	"sync"
	"time"
)

type TCPProxy struct {
	dialTimeout time.Duration
	maxConn int
	connSem chan struct{}
	wg sync.WaitGroup
}

func NewTCPProxy(maxConn int) *TCPProxy{
	return & TCPProxy{
		dialTimeout: time.Second * 10,
		maxConn: maxConn,
		connSem: make(chan struct{}, maxConn),
	}
}


func (p *TCPProxy) HandleConnection(client net.Conn) {
	select {
	case p.connSem <- struct{}{}:
		defer func(){ <-p.connSem }()
	default:
		client.Close()
		return
	}

	defer client.Close()
	p.wg.Add(1)
	defer p.wg.Done()

	buf := make([]byte, 256)
	n, err := client.Read(buf)
	if err != nil {
		return
	}

	destAddr := string(buf[:n])
	destAddr = destAddr[:len(destAddr) - 1]

	upstream, err := net.DialTimeout("tcp", destAddr, p.dialTimeout)
	if err != nil {
		return
	}

	defer upstream.Close()

	chanErr := make(chan error, 2)

	// client to upstream
	go func() {
		_, err := io.Copy(upstream, client)
		if tcpConn, ok := upstream.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
		chanErr <- err
	}()

	go func() {
		_, err := io.Copy(client, upstream)
		if tcpConn, ok := client.(*net.TCPConn); ok {
			tcpConn.CloseWrite()
		}
		chanErr <- err
	}()

	<-chanErr
} 

func (p *TCPProxy) HandleShutdown() {
	p.wg.Wait()
}