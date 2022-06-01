package plug

import (
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"syscall"
)

type SocketPlug interface {
	Connect(ip string, port int) (syscall.Handle, error)
	Send(msg string) error
	Receive(len int) (string, error)
}

func init() {
	var wsadata syscall.WSAData
	if err := syscall.WSAStartup(MAKEWORD(2, 2), &wsadata); err != nil {
		log.Printf("Startup error, [%s]\n", err)
		return
	}
}

func MAKEWORD(low, high uint8) uint32 {
	var ret = uint16(high)<<8 + uint16(low)
	return uint32(ret)
}

const DISCOVER_HEX = "21310020ffffffffffffffffffffffffffffffffffffffffffffffffffffffff"

type Socket struct {
	fd syscall.Handle
	sa syscall.Sockaddr
}

func NewSocket(ip string, port int) (*Socket, error) {
	plugSocket := Socket{}
	return plugSocket.Connect(ip, port)
}

// Connect setup Socket Handle
func (s *Socket) Connect(ip string, port int) (*Socket, error) {
	// chuangmi_plug connect use udp,set socket type to SOCK_DGRAM
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_IP)
	if err != nil {
		return nil, err
	}

	if err = syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1); err != nil {
		return nil, err
	}

	s.sa = &syscall.SockaddrInet4{
		Port: port,
		Addr: inetAddr(ip),
	}

	if err = syscall.Connect(fd, s.sa); err != nil {
		return nil, err
	}
	s.fd = fd
	return s, nil
}

func send(fd syscall.Handle, buf syscall.WSABuf) error {
	//var buf syscall.WSABuf
	var written uint32
	//buf.Buf, _ = syscall.BytePtrFromString(msg)
	//buf.Len = uint32(len(msg))
	err := syscall.WSASend(fd, &buf, 1, &written, 0, nil, nil)
	if err != nil {
		log.Printf("write error [%s]\n", err)
		return err
	}
	return nil
}

// Discover the devices at given network
func (s Socket) Discover() error {
	var (
		d   []byte
		err error
	)
	if d, err = hex.DecodeString(DISCOVER_HEX); err != nil {
		return err
	}
	if err := syscall.Sendto(s.fd, d, 0, s.sa); err != nil {
		return err
	}
	return nil
}

// inetAddr is used to convert ip string to ip byte
func inetAddr(ipaddr string) [4]byte {
	var (
		ips = strings.Split(ipaddr, ".")
		ip  [4]uint64
		ret [4]byte
	)
	for i := 0; i < 4; i++ {
		ip[i], _ = strconv.ParseUint(ips[i], 10, 8)
	}
	for i := 0; i < 4; i++ {
		ret[i] = byte(ip[i])
	}
	return ret
}
