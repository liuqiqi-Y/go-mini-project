package main

func main() {
	servers := newServers(":9998", ":9999")
	for _, s := range servers {
		go s.run()
	}
	p := newProxy(":9000", ":9999")
	go p.run()
	for _, c := range newClient(":9000", 3) {
		go c.run()
	}
	select {}
}
