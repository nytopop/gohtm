package main

import "github.com/nytopop/gohtm/net"

func main() {
	for i := 512; i < 560; i++ {
		net.NewNetwork(i)
	}
}
