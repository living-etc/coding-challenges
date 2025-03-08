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

		if server.exec != "" {
			go server.handleExecConnection(conn, server.exec)
		} else {
			go server.handleStdConnection(conn)
		}
	}
}

func (server *TcpServer) handleExecConnection(conn net.Conn, execCmd string) {
	defer conn.Close()

	cmd := exec.Command(execCmd)

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
}

func (server *TcpServer) handleStdConnection(conn net.Conn) {
	defer conn.Close()

	go io.Copy(os.Stdout, conn)

	io.Copy(conn, os.Stdin)
}
