package tracerpki

import (
	"testing"
)

func TestGetDestinationAddresses(t *testing.T) {
	tests := []struct {
		dest string
	}{
		{
			dest: "google.com",
		},
		{
			dest: "bbc.co.uk",
		},
		{
			dest: "8.8.8.8",
		},
		{
			dest: "2a04:4e42:400::81",
		},
		{
			dest: "lol",
		},
	}
	for _, tc := range tests {
		addrs, err := getDestinationAddresses(tc.dest)
		t.Errorf("%v - %v", addrs, err)
	}
}
