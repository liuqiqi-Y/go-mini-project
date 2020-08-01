package main

import (
	"fmt"
	"net"
)

type server struct {
	listener net.Listener
}

func newServers(addresses ...string) []*server {
	s := make([]*server, 0)
	for _, a := range addresses {
		lst, err := net.Listen("tcp", a)
		if err != nil {
			fmt.Printf("failed to open a listener: %s\n", err.Error())
			return nil
		}
		s = append(s, &server{listener: lst})
	}
	return s
}

func (s *server) run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("failed to listen: %s\n", err.Error())
			return
		}
		go s.handler(conn)
	}
}
func (s *server) handler(conn net.Conn) {
	for {
		msg := make([]byte, 1024)
		_, err := conn.Read(msg)
		if err != nil {
			fmt.Printf("failed to read data from client: %s\n", err.Error())
			return
		}
		_, err = conn.Write([]byte(fmt.Sprintf("message from %s\n", s.listener.Addr().String())))
		if err != nil {
			fmt.Printf("failed to write data to client: %s\n", err.Error())
			return
		}
	}
}
