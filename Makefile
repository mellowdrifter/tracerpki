build:
	gofumpt -w *.go
	go build -o tracerpki

cover:
	go test -cover ./...

race:
	go test -race