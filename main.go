package main

import (
	"fmt"
	"klipper_power/pkg/plug"
	"syscall"
)

const (
	IP       string = "127.0.0.1"
	PORT     int    = 54321
	RECV_LEN        = 1024
	TOKEN    string = "ac03b24f88acfb0937942bddc46be066"
)

func main() {
	//udpAddr, err := net.ResolveUDPAddr("udp4", IP+":"+strconv.Itoa(PORT))
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//conn, err := net.DialUDP("udp4", nil, udpAddr)
	//defer conn.Close()

	token := [16]byte{}
	bToken, _ := syscall.ByteSliceFromString(TOKEN)
	copy(token[:], bToken[:])

	req := plug.DiscoverRequest(token)

	fmt.Println(req)
	//Len := unsafe.Sizeof(req)
	//testBytes := &SliceMock{
	//	addr: uintptr(unsafe.Pointer(&req)),
	//	cap:  int(Len),
	//	len:  int(Len),
	//}
	//data := *(*[]byte)(unsafe.Pointer(testBytes))
	//_, err = conn.Write(req)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//var buf [1024]byte
	//_, err = conn.Read(buf[0:])
	//Struct := *(**Request)(unsafe.Pointer(&data))
	//fmt.Println(Struct)
}

type DID struct {
	did   string
	siid  uint8
	piid  uint8
	value bool
}
