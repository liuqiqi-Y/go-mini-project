package client

import (
	"fmt"
	"net"
)

type Client struct {
	cli net.Conn
}

func NewClients(proxy string, n int) []*Client {
	clients := make([]*Client, n)
	for i := 0; i < n; i++ {
		conn, err := net.Dial("tcp", proxy)
		if err != nil {
			panic(err)
		}
		clients[i] = &Client{conn}
	}
	return clients
}
func (c *Client) Send() {
	_, err := c.cli.Write([]byte(fmt.Sprintf("$$$message from %s", c.cli.LocalAddr().String())))
	if err != nil {
		panic(err)
	}
	input := make([]byte, 1024)
	bytes, err := c.cli.Read(input)
	if bytes > 0 {
		fmt.Println(string(input))
	}
	if err != nil {
		panic(err)
	}
}
