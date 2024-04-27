package main

import (
	"io"
	"log"
	"net"
	"os"
)

type TcpServer struct {
	port string
}

func NewTcpServer(port string) *TcpServer {
	return &TcpServer{
		port: port,
	}
}

func (server *TcpServer) Start() {
	l, err := net.Listen("tcp", ":"+server.port)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go server.handleInbound(conn)
		go server.handleOutbound(conn)
	}
}

func (server *TcpServer) handleInbound(c net.Conn) {
	_, err := io.Copy(os.Stdout, c)
	if err != nil {
		log.Println(err)
	}
	c.Close()
}

func (server *TcpServer) handleOutbound(c net.Conn) {
	_, err := io.Copy(c, os.Stdin)
	if err != nil {
		log.Println(err)
	}
	c.Close()
}
