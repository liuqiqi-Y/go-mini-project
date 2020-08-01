package main

import (
	"fmt"
	"io"
	"net"
)

type proxyAgent struct {
	lst            net.Listener
	backendAddress string
}

func newProxy(agent, backend string) *proxyAgent {
	conn, err := net.Listen("tcp", agent)
	if err != nil {
		fmt.Printf("faile to open a listener: %s\n", err.Error())
		return nil
	}
	return &proxyAgent{
		lst:            conn,
		backendAddress: backend,
	}
}
func (p *proxyAgent) agent(front, back net.Conn) {
	go io.Copy(back, front)
	go io.Copy(front, back)
}
func (p *proxyAgent) run() {
	for {
		front, err := p.lst.Accept()
		if err != nil {
			fmt.Printf("failed to listen: %s\n", err.Error())
			return
		}
		back, err := net.Dial("tcp", p.backendAddress)
		if err != nil {
			fmt.Printf("failed to dial back server: %s\n", err.Error())
		}
		go p.agent(front, back)
	}
}
