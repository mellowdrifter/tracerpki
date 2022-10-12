module github.com/mellowdrifter/tracerpki

go 1.19

replace github.com/mellowdrifter/bgp_infrastructure/proto/glass => ../bgp_infrastructure/proto/glass

replace github.com/mellowdrifter/go-bgpstuff.net => ../go-bgpstuff.net

require (
	github.com/mellowdrifter/go-bgpstuff.net v0.0.0-20220507215736-e57e864fa24b
	github.com/spf13/cobra v1.6.0
)

require (
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/mellowdrifter/bogons v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/time v0.0.0-20220609170525-579cf78fd858 // indirect
)
