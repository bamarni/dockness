VERSION ?= v2.0.1

clean:
	rm -rf build

build: clean
	mkdir -p build/releases
	for os in "darwin" "linux" ; do \
		mkdir build/$$os; \
		GOOS=$$os GOARCH=amd64 go build -o build/$$os/dockness ./dockness.go; \
		tar -cvzf build/releases/dockness-$$os-x64.tar.gz -C build/$$os dockness; \
	done

release:
	git tag $(VERSION)
	git push origin $(VERSION)
	github-release release --user bamarni --repo dockness --tag $(VERSION) --name "$(VERSION)"
	for os in "darwin" "linux" ; do \
		github-release upload --user bamarni --repo dockness --tag $(VERSION) \
			--name "dockness-$$os-x64.tar.gz" \
			--file build/releases/dockness-$$os-x64.tar.gz; \
	done
