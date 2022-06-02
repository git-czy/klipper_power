package main

import (
	"fmt"
	"klipper_power/pkg/http"
	"klipper_power/pkg/plug"
)

const (
	IP          string = "127.0.0.1"
	PORT        int    = 54321
	TOKEN       string = "ac03b24f88acfb0937942bddc46be066"
	KLIPPER_URI string = "http://106.55.18.128:8001/v1/deploy/git_auth"
	DURING      int    = 5
)

func main() {
	//conn, err := plug.NewSocket(IP, PORT, TOKEN)
	//if err != nil {
	//	return
	//}
	//defer conn.Close()
	//
	//if err = conn.Discover(); err != nil {
	//	fmt.Printf("discover failed: %s", err.Error())
	//	return
	//}
	//
	//err = conn.PowerOff()
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	conn, err := plug.NewSocket(IP, PORT, TOKEN)
	if err != nil {
		panic(err)
	}

	go http.HandleKlipper(KLIPPER_URI, conn, DURING)
	go http.StartServer(conn)

}
