# Run locally with `docker run --rm -it -v "$(pwd)/config:/app/config" dnstest:latest run --config /app/config/config.test.yaml`
FROM golang:1.22-alpine as builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/dnsTest .

FROM scratch
WORKDIR /bin
COPY --from=builder /src/bin/dnsTest .
ENTRYPOINT ["/bin/dnsTest"]
