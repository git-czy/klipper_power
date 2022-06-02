package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
)

const (
	SERVER_IP       = "127.0.0.1"
	SERVER_PORT     = 54321
	SERVER_RECV_LEN = 1024
)

func main() {
	fmt.Println("start udp http")
	address := SERVER_IP + ":" + strconv.Itoa(SERVER_PORT)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer conn.Close()

	for {
		// Here must use make and give the lenth of buffer
		data := make([]byte, SERVER_RECV_LEN)
		_, rAddr, err := conn.ReadFromUDP(data)
		if err != nil {
			fmt.Println(err)
			continue
		}

		strData := string(data)
		fmt.Println("Received:", []byte(strData))
		fmt.Println("len:", len([]byte(strData)))

		_, err = conn.WriteToUDP([]byte(strData), rAddr)
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Send:", []byte(strData))
	}

}
