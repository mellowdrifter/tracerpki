build:
	gofumpt -w *.go
	go build -o tracerpki
	sudo chown root ./tracerpki
	sudo chmod u+s ./tracerpki

cover:
	go test -cover ./...

race:
	go test -race