package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Config struct {
	Listen *ListenConfig
	Scan   *ScanConfig
}

type ListenConfig struct {
	Port string
	Udp  bool
	Exec string
}

type ScanConfig struct {
	Host  string
	Ports []int
}

func ParseConfig(args []string) Config {
	flagSet := flag.NewFlagSet("flags", flag.ContinueOnError)

	portFlag := flagSet.String("p", "", "Port number")
	listenFlag := flagSet.Bool("l", false, "Run in listen mode")
	scanFlag := flagSet.Bool("z", false, "Scan for an open port")
	udpFlag := flagSet.Bool("u", false, "Listen for UDP connections")
	execFlag := flagSet.String("e", "", "Execute the specified command")

	err := flagSet.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	positionalFlags := flagSet.Args()

	var config Config
	if *listenFlag && *scanFlag {
		log.Fatal("Cannot specify both -l and -z")
		os.Exit(1)
	}

	if *listenFlag {
		config.Listen = &ListenConfig{
			Port: *portFlag,
			Udp:  *udpFlag,
			Exec: *execFlag,
		}
	}

	if *scanFlag {
		ports, err := parsePorts(positionalFlags[1])
		if err != nil {
			log.Fatal(err)
		}

		config.Scan = &ScanConfig{
			Host:  positionalFlags[0],
			Ports: ports,
		}
	}

	return config
}

func parsePorts(portInput string) ([]int, error) {
	portPattern := `^(\d+)(?:-(\d+))?$`
	re := regexp.MustCompile(portPattern)
	matches := re.FindStringSubmatch(portInput)

	if matches == nil {
		return nil, fmt.Errorf("invalid port format")
	}

	startPort, err := strconv.Atoi(matches[1])
	if err != nil {
		return nil, fmt.Errorf("invalid start port")
	}

	if matches[2] == "" {
		return []int{startPort}, nil
	}

	endPort, err := strconv.Atoi(matches[2])
	if err != nil {
		return nil, fmt.Errorf("invalid end port")
	}

	if startPort > endPort {
		return nil, fmt.Errorf("start port cannot be greater than end port")
	}

	var ports []int
	for p := startPort; p <= endPort; p++ {
		ports = append(ports, p)
	}

	return ports, nil
}
