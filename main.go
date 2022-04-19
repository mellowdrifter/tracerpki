package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/netip"
	"syscall"
	"time"

	gpb "github.com/mellowdrifter/bgp_infrastructure/proto/glass"
	"google.golang.org/grpc"
)

const (
	dialDeadline = 2
	reqPerMinute = 30
	startPort    = 33434
	endPort      = 33534 // do i really need this if maxHops == 30 ?
	maxHops      = 30
	timeoutMS    = 3000
)

var (
	payload         = []byte("@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_")
	defaultLocation = "google.com"
)

func main() {
	// Create the UDP sender and the ICMP receiver
	sender, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, syscall.IPPROTO_UDP)
	if err != nil {
		panic(err)
	}
	receiver, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		panic(err)
	}

	defer syscall.Close(receiver)
	defer syscall.Close(sender)

	// get glass server

	glass, err := getGRPC("127.0.0.1:7181")
	if err != nil {
		panic(err)
	}

	// Initial ttl and timeouts
	ttl := 1
	timeout := syscall.NsecToTimeval(1000000 * timeoutMS)

	// find the local address
	socketAddr, err := getSocketAddress()
	if err != nil {
		panic(err)
	}

	// where are we trying to get to?
	destAddr, err := getDestinationAddress(defaultLocation)
	if err != nil {
		panic(err)
	}

	fmt.Printf("traceroute to %s (%s), %d hops max, using %s as source\n", defaultLocation, destAddr, maxHops, socketAddr)

	ctx := context.TODO()

	for {
		// set the ttl and timeouts on the socket
		syscall.SetsockoptInt(sender, 0x0, syscall.IP_TTL, ttl)
		syscall.SetsockoptTimeval(receiver, syscall.SOL_SOCKET, syscall.SO_RCVTIMEO, &timeout)

		syscall.Bind(receiver, &syscall.SockaddrInet4{Port: startPort, Addr: socketAddr.As4()})
		syscall.Sendto(sender, payload, 32, &syscall.SockaddrInet4{Port: startPort, Addr: destAddr.As4()})
		// Increment the source port each time. We don't have to, but traceroute seems to do this
		// startPort++

		// now for some messages?
		// 52 = packetsize
		p := make([]byte, 52)
		// n = packet size
		_, from, err := syscall.Recvfrom(receiver, p, 0)
		// err can be returned if timeout is reached
		if err != nil {
			ttl++
			continue
		}
		ip := from.(*syscall.SockaddrInet4).Addr
		ipString := fmt.Sprintf("%v.%v.%v.%v", ip[0], ip[1], ip[2], ip[3])
		// not all addresses will have reverse names. Therefore no need to check for an error. If there is an error, just don't bother displaying a hostname.
		host, _ := net.LookupAddr(ipString)
		origin, err := glass.Origin(ctx, &gpb.OriginRequest{IpAddress: &gpb.IpAddress{Address: ipString, Mask: 32}})
		if err != nil {
			panic(err)
		}
		asname, err := glass.Asname(ctx, &gpb.AsnameRequest{AsNumber: origin.GetOriginAsn()})
		if err != nil {
			panic(err)
		}
		// TODO: I need the full route, not just a single IP
		rpki, err := glass.Roa(ctx, &gpb.RoaRequest{IpAddress: &gpb.IpAddress{Address: ipString, Mask: 32}})
		if err != nil {
			panic(err)
		}
		fmt.Printf("%d\t%s\t%s\t%d\t%s\tRPKI STATUS: %s\n", ttl, host, ipString, origin.GetOriginAsn(), asname.GetAsName(), rpki.GetStatus().String())

		if ttl > maxHops {
			fmt.Println("maximum hops reached")
			break
		}
		if ip == destAddr.As4() {
			break
		}

		ttl++
	}
}

func getSocketAddress() (netip.Addr, error) {
	iAddrs, err := net.InterfaceAddrs()
	if err != nil {
		return netip.Addr{}, err
	}

	for _, a := range iAddrs {
		addr, err := netip.ParsePrefix(a.String())
		if err != nil {
			return netip.Addr{}, errors.New("Unable to choose a source IP address")
		}
		if !addr.Addr().IsLoopback() {
			return addr.Addr(), nil
		}
	}
	return netip.Addr{}, nil
}

func getDestinationAddress(dest string) (netip.Addr, error) {
	addrs, err := net.LookupHost(dest)
	if err != nil {
		return netip.Addr{}, err
	}
	// Just grab the first address...
	// addr := addrs[0]
	for _, addr := range addrs {
		t, _ := netip.ParseAddr(addr)
		if t.Is4() {
			return t, nil
		}
	}
	return netip.Addr{}, errors.New("no ipv4 address found")

	// return netip.ParseAddr(addr)
}

// getGRPC will dial however many looking glass servers and return
// a slice of clients ready to go.
func getGRPC(srv string) (gpb.LookingGlassClient, error) {
	conn, err := grpc.Dial(srv,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(dialDeadline*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("Unable to dial gRPC server %s: %v", srv, err)
	}
	client := gpb.NewLookingGlassClient(conn)
	return client, nil
}
