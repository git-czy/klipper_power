package plug

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"strconv"
	"syscall"
	"time"
)

type SocketPlug interface {
	Send(msg string) error
	Receive(len int) (string, error)
	Discover() error
}

type Socket struct {
	conn     *net.UDPConn
	token    []byte
	discover *discover
}

type discover struct {
	deviceId uint32
	stamp    uint32
}

const MaxBufferSize = 2048

func NewSocket(ip string, port int, token string) (*Socket, error) {
	var (
		udpAddr *net.UDPAddr
		conn    *net.UDPConn
		t       []byte
		err     error
	)
	if udpAddr, err = net.ResolveUDPAddr("udp4", ip+":"+strconv.Itoa(port)); err != nil {
		return nil, err
	}
	if conn, err = net.DialUDP("udp4", nil, udpAddr); err != nil {
		return nil, err
	}
	if t, err = syscall.ByteSliceFromString(token); err != nil {
		return nil, err
	}
	return &Socket{conn, t, nil}, nil
}

// Discover the devices at given network
func (s *Socket) Discover() error {
	request := NewPowerRequest(s.token, "hello")
	if _, err := s.conn.Write(request); err != nil {
		return err
	}
	response, err := readResponse(s.conn)
	if err != nil {
		return err
	}
	s.discover, err = parseDiscoverRes(response)
	if err != nil {
		return err
	}
	return nil
}

func parseDiscoverRes(res []byte) (*discover, error) {
	if len(res) != 32 {
		return nil, fmt.Errorf("expected 32 bytes, got %d", len(res))
	}
	return &discover{
		deviceId: binary.BigEndian.Uint32(res[8:]),
		stamp:    binary.BigEndian.Uint32(res[12:]),
	}, nil
}

func (s Socket) PowerOff() error {
	request := NewPowerRequest(s.token, "off")
	if _, err := s.conn.Write(request); err != nil {
		return err
	}
	response, err := readResponse(s.conn)
	if err != nil {
		return err
	}
	fmt.Println(response)
	if len(response) < 32 {
		return fmt.Errorf("lost data while power off,at least 32 bytes,but got %d", len(response))
	}
	return nil
}

func (s Socket) PowerOn() error {
	request := NewPowerRequest(s.token, "on")
	if _, err := s.conn.Write(request); err != nil {
		return err
	}
	response, err := readResponse(s.conn)
	if err != nil {
		return err
	}
	if len(response) < 32 {
		log.Println("Lost data while power on!")
	}
	return nil
}

func (s Socket) Close() {
	_ = s.conn.Close()
}

func readResponse(conn net.Conn) ([]byte, error) {
	return readTimeout(conn)
}

func readTimeout(conn net.Conn) ([]byte, error) {
	resultChan := make(chan int)
	errChan := make(chan error, 1)
	buffer := make([]byte, MaxBufferSize)
	go func() {
		err := conn.SetReadDeadline(deadline())
		if err != nil {
			errChan <- err
			return
		}
		n, err := conn.Read(buffer)
		if err != nil {
			errChan <- err
			return
		}
		resultChan <- n
	}()

	select {
	case result := <-resultChan:
		return buffer[:result], nil
	case err := <-errChan:
		return nil, err
	}
}

var timeout = 2 * time.Second

func deadline() time.Time {
	return time.Now().Add(timeout)
}
