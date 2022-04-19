package main

import (
	"net/netip"

	gpb "github.com/mellowdrifter/bgp_infrastructure/proto/glass"
)

type hop struct {
	hop        int
	asNumber   int
	address    netip.Addr
	asName     string
	rDNS       string
	rpkiStatus gpb.RoaResponse
}
