package main

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
)

type TcpServer struct {
	port string
	exec string
}

func NewTcpServer(port, exec string) *TcpServer {
	return &TcpServer{
		port: port,
		exec: exec,
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

func (server *TcpServer) handleInbound(conn net.Conn) {
	defer conn.Close()

	if server.exec != "" {
		cmd := exec.Command(server.exec)

		stdin, _ := cmd.StdinPipe()
		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		cmd.Start()

		done := make(chan bool, 2)

		go func() {
			io.Copy(stdin, conn)
			stdin.Close()
			done <- true
		}()

		go func() {
			io.Copy(conn, stdout)
			done <- true
		}()

		go func() {
			io.Copy(conn, stderr)
			done <- true
		}()

		<-done
		<-done

		cmd.Wait()
	} else {
		io.Copy(os.Stdout, conn)
	}

	conn.Close()
}

func (server *TcpServer) handleOutbound(c net.Conn) {
	_, err := io.Copy(c, os.Stdin)
	if err != nil {
		log.Println(err)
	}
	c.Close()
}
