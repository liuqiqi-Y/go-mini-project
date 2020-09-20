package main

import (
	"go-mini-project/proxy/agent"
	"go-mini-project/proxy/client"
	"go-mini-project/proxy/server"
)

func main() {
	servers := server.NewServers(":9997", ":9998", ":9999")
	for _, s := range servers {
		go s.Run()
	}
	p := agent.NewAgent(":9000", ":9997", ":9998", ":9999")
	go p.Run()
	for _, c := range client.NewClients(":9000", 3) {
		go c.Send()
	}
	select {}
}
