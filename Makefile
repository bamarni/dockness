VERSION ?= v2.0.2
export MACHINE_STORAGE_PATH=$(PWD)/test

.PHONY: clean test build release

vendor:
	govendor sync

build: vendor
	rm -rf build
	mkdir -p build/releases
	for os in "darwin" "linux" ; do \
		mkdir build/$$os; \
		GOOS=$$os GOARCH=amd64 go build -o build/$$os/dockness ./dockness.go; \
		tar -cvzf build/releases/dockness-$$os-x64.tar.gz -C build/$$os dockness; \
	done

test: vendor
	docker-machine create -d generic --generic-ip-address 192.0.2.0 test >/dev/null 2>&1&
	sleep 3
	go test -race
	go test -run XXX -bench .

release: test build
	git tag $(VERSION)
	git push origin $(VERSION)
	github-release release --user bamarni --repo dockness --tag $(VERSION) --name "$(VERSION)"
	for os in "darwin" "linux" ; do \
		github-release upload --user bamarni --repo dockness --tag $(VERSION) \
			--name "dockness-$$os-x64.tar.gz" \
			--file build/releases/dockness-$$os-x64.tar.gz; \
	done

clean:
	rm -rf build test vendor
