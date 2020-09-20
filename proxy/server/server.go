package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
}

func NewServers(addresses ...string) []*Server {
	servers := make([]*Server, len(addresses))
	for k, v := range addresses {
		listener, err := net.Listen("tcp", v)
		if err != nil {
			panic(err)
		}
		servers[k] = &Server{listener}
	}
	return servers
}

func (s *Server) Run() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			panic(err)
		}
		go handler(conn)
	}
}

func handler(conn net.Conn) {
	for {
		input := make([]byte, 1024)
		bytes, err := conn.Read(input)
		if bytes > 0 {
			_, err := conn.Write([]byte(fmt.Sprintf("$$$response for %s", string(input))))
			if err != nil {
				panic(err)
			}
		}
		if err != nil {
			panic(err)
		}
	}
}
