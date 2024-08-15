FROM golang:1.22.6-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o bin/sherlock .

FROM scratch
WORKDIR /bin
COPY --from=builder /src/bin/sherlock .
ENTRYPOINT ["/bin/sherlock"]
