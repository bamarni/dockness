dev:
	tmuxinator start go

release:
	mkdir -p build
	GOOS=darwin GOARCH=amd64 go build -o build/docker-machine-dns-darwin-x64 ./docker-machine-dns.go
	GOOS=linux GOARCH=amd64 go build -o build/docker-machine-dns-linux-x64 ./docker-machine-dns.go

