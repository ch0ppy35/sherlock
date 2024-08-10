.PHONY : build
build:
	go build -o bin/sherlock .

.PHONY: test
test:
	go test ./... -v

.PHONY: clean
clean:
	rm -rf bin/