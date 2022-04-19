module github.com/mellowdrifter/tracerpki

go 1.18

replace github.com/mellowdrifter/bgp_infrastructure/proto/glass => ../bgp_infrastructure/proto/glass

require github.com/mellowdrifter/bgp_infrastructure/proto/glass v0.0.0-00010101000000-000000000000

require (
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20200822124328-c89045814202 // indirect
	golang.org/x/sys v0.0.0-20200323222414-85ca7c5b95cd // indirect
	golang.org/x/text v0.3.0 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	google.golang.org/grpc v1.43.0 // indirect
	google.golang.org/protobuf v1.26.0 // indirect
)
