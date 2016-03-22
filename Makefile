clean:
	rm -rf build

build: clean
	mkdir -p build/darwin build/linux
	GOOS=darwin GOARCH=amd64 go build -o build/darwin/dockness ./dockness.go
	GOOS=linux GOARCH=amd64 go build -o build/linux/dockness ./dockness.go

release: build
	mkdir -p build/releases
	tar -cvzf build/releases/dockness-darwin-x64.tar.gz -C build/darwin dockness
	tar -cvzf build/releases/dockness-linux-x64.tar.gz -C build/linux dockness
