package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const maxBufferSize = 4096

type UdpServer struct {
	port        string
	remoteAddr  net.Addr
	remoteAddrC chan net.Addr
	errorChan   chan error
}

func NewUdpServer(port string) *UdpServer {
	return &UdpServer{
		port:        port,
		remoteAddrC: make(chan net.Addr, 1),
		errorChan:   make(chan error, 1),
	}
}

func (server *UdpServer) Start() {
	l, err := net.ListenPacket("udp", ":"+server.port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	go server.handleInbound(l)
	go server.handleOutbound(l)

	err = <-server.errorChan
	if err != nil {
		log.Fatal(err)
	}
}

func (server *UdpServer) handleInbound(l net.PacketConn) {
	buffer := make([]byte, maxBufferSize)

	for {
		n, addr, err := l.ReadFrom(buffer)
		if err != nil {
			server.errorChan <- err
			return
		}

		if server.remoteAddr == nil {
			server.remoteAddr = addr
			server.remoteAddrC <- addr
		}

		fmt.Println(string(buffer[:n]))
	}
}

func (server *UdpServer) handleOutbound(l net.PacketConn) {
	scanner := bufio.NewScanner(os.Stdin)

	var currentAddr net.Addr
	select {
	case addr := <-server.remoteAddrC:
		currentAddr = addr
	case err := <-server.errorChan:
		server.errorChan <- err
		return
	}

	for scanner.Scan() {
		msg := scanner.Text() + "\n"

		_, err := l.WriteTo([]byte(msg), currentAddr)
		if err != nil {
			server.errorChan <- err
			return
		}
	}
	if err := scanner.Err(); err != nil {
		server.errorChan <- err
	}
}
