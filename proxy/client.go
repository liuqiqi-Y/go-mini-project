package main

import (
	"fmt"
	"net"
)

type client struct {
	conn net.Conn
}

func newClient(proxy string, n int) []*client {
	c := make([]*client, 0)
	for i := 0; i < n; i++ {
		conn, err := net.Dial("tcp", proxy)
		if err != nil {
			fmt.Printf("failed to connect proxy server: %s\n", err.Error())
			return nil
		}
		c = append(c, &client{conn: conn})
	}
	return c
}
func (c *client) run() {
	_, err := c.conn.Write([]byte(fmt.Sprintf("message from client: %s\n", c.conn.LocalAddr().String())))
	if err != nil {
		fmt.Printf("faile to write data to proxy: %s\n", err.Error())
		return
	}
	msg := make([]byte, 1024)
	_, err = c.conn.Read(msg)
	if err != nil {
		fmt.Printf("faile to write data to proxy: %s\n", err.Error())
		return
	}
	fmt.Println(string(msg))
}
