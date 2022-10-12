package main

import (
	"fmt"
	"net"
)

func main() {
	addrs, _ := net.LookupHost("www.google.com")
	fmt.Println(addrs)
}
