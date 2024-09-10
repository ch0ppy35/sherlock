FROM golang:1.23.1-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY main.go main.go
COPY cmd/ cmd/
COPY internal/ internal/

RUN BUILD_DATE=$(date -u +'%Y-%m-%dT%H:%M:%SZ') && \
    BUILD_ARCH=$(go env GOARCH) && \
    GO_VERSION=$(go env GOVERSION) && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w \
    -X github.com/ch0ppy35/sherlock/cmd.version=LOCALDOCKERBUILD \
    -X github.com/ch0ppy35/sherlock/cmd.commit=HEAD \
    -X github.com/ch0ppy35/sherlock/cmd.date=${BUILD_DATE} \
    -X github.com/ch0ppy35/sherlock/cmd.arch=${BUILD_ARCH} \
    -X github.com/ch0ppy35/sherlock/cmd.goversion=${GO_VERSION}" \
    -o bin/sherlock .

FROM scratch
COPY --from=builder /src/bin/sherlock /sherlock
ENTRYPOINT ["/sherlock"]