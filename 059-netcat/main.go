package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	config := ParseConfig(os.Args[1:])

	if config.Listen != nil {
		listenConfig := config.Listen

		if listenConfig.Udp {
			server := NewUdpServer(listenConfig.Port)

			server.Start()
		} else {
			server := NewTcpServer(listenConfig.Port, listenConfig.Exec)

			server.Start()
		}
	} else if config.Scan != nil {
		scanConfig := config.Scan

		host := scanConfig.Host
		ports := scanConfig.Ports

		for _, port := range ports {
			conn, err := net.Dial("tcp", fmt.Sprintf("%v:%v", host, port))

			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Printf("Succeeded in connecting to %v on port %v", host, port)

				conn.Close()
			}
		}
	} else {
		fmt.Println("Usage: ccnc [-luz] [-i interval] [-p source_port] [hostname] [port[s]]")
	}
}
