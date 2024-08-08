.PHONY : build
build:
	go build -o bin/dnsTest .

.PHONY: test
test:
	go test ./... -v

.PHONY: clean
clean:
	rm -rf bin/