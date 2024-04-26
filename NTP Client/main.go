package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"net"
	"os"
	"time"
)

const (
	jan1970     = 2208988800
	defaultPort = "123"
)

type NTPPacketHeader struct {
	LiVnMode           uint8
	Stratum            uint8
	Poll               uint8
	Precision          int8
	RootDelay          uint32
	RootDispersion     uint32
	ReferenceID        uint32
	ReferenceTimestamp uint64
	OriginTimestamp    uint64
	ReceiveTimestamp   uint64
	TransmitTimestamp  uint64
}

var (
	url  string
	port string
)

func init() {
	flag.StringVar(&port, "p", defaultPort, "Port of the NTP server")
}

func printUsage(programName string) {
	fmt.Printf("Usage: %s <URL>\n\n", programName)
	fmt.Println("A simple NTP client\n")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		printUsage(os.Args[0])
		os.Exit(1)
	}

	url = args[0]

	res, err := net.LookupIP(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "getaddrinfo error: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   res[0],
		Port: 123,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "socket error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	header := NTPPacketHeader{
		LiVnMode: 0x23,
	}

	packetBytes := make([]byte, 48)
	packetBytes[0] = header.LiVnMode

	_, err = conn.Write(packetBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "sendto error: %v\n", err)
		os.Exit(1)
	}

	responseBytes := make([]byte, 48)
	_, _, err = conn.ReadFromUDP(responseBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "recvfrom error: %v\n", err)
		os.Exit(1)
	}

	transmitTimestamp := binary.BigEndian.Uint32(responseBytes[40:44])
	currentTime := int64(transmitTimestamp - jan1970)

	fmt.Printf("The current time is %s\n", time.Unix(currentTime, 0))
}