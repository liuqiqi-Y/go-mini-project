package agent

import (
	"io"
	"net"
	"sync"
)

type Agent struct {
	listener net.Listener
	backends []string
}

func NewAgent(agent string, backends ...string) *Agent {
	listener, err := net.Listen("tcp", agent)
	if err != nil {
		panic(err)
	}
	bks := make([]string, len(backends))
	for i, j := range backends {
		bks[i] = j
	}
	return &Agent{
		listener: listener,
		backends: bks,
	}
}

func (a *Agent) Run() {
	i := 0
	for {
		conn1, err := a.listener.Accept()
		if err != nil {
			panic(err)
		}
		conn2, err := net.Dial("tcp", a.backends[i%len(a.backends)])
		if err != nil {
			panic(err)
		}
		go handler(conn1, conn2)
		i++
	}
}
func handler(a net.Conn, b net.Conn) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		_, err := io.Copy(b, a)
		if err != nil {
			panic(err)
		}
	}()
	go func() {
		defer wg.Done()
		_, err := io.Copy(a, b)
		if err != nil {
			panic(err)
		}
	}()
	wg.Wait()
}
