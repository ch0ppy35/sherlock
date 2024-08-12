.PHONY: build
build:
	go build -ldflags "-X github.com/ch0ppy35/sherlock/cmd.version=LOCALBUILD \
	-X github.com/ch0ppy35/sherlock/cmd.commit=main \
	-X github.com/ch0ppy35/sherlock/cmd.date=$$(date -u +'%Y%m%d%H%M')" \
	-o bin/sherlock .

.PHONY: clean
clean:
	rm -rf bin/
	rm -rf dist/

.PHONY: docker-build
docker-build:
	docker build -t ghcr.io/ch0ppy35/sherlock:latest .

.PHONY: test
test:
	go test -v ./...
