# Sherlock

Simple Infrastructure Sanity Checks

## Overview

Sherlock is a cli tool designed for simple infrastructure sanity checks. Currently it allows you to perform DNS record tests using a YAML configuration file or individual parameters, verifying records such as A, AAAA, CNAME, MX, TXT, and NS.

In addition to DNS tests more capabilities will be added in the future, including SFTP testing. Sherlock is intended for use within containers for CI/CD pipelines or cronjobs in Kubernetes. Binaries are available in the GitHub release, or you can build the binary locally with the `make` command.

See `sherlock -h` for more info

## DNS Testing Usage

### Config File

Create a YAML configuration file specifying the DNS tests you want to perform. Below is an example configuration:

```yaml
dnsServer: "8.8.8.8"
tests:
  - host: sftp.foobar.com
    expectedValues: ["10.0.0.10", "10.1.0.10", "10.3.0.10"]
    testType: a
  - host: sftp.foobar.com
    expectedValues: ["s-12345678900f0000a.server.transfer.us-east-1.amazonaws.com."]
    testType: cname
  - host: grafana.foobar.com
    expectedValues: ["10.0.0.100"]
    testType: a
```

### Running Tests

```bash
sherlock dns run --config path/to/config.yaml

# or without a config file

sherlock dns test --server 1.1.1.1 --host prom.example.com --expected "10.0.0.1" --type a
```

### Docker

Alternatively, you can run Sherlock inside a Docker container

```bash
docker run --rm -it \
  -v "$(pwd)/config:/app/config" \
  ghcr.io/ch0ppy35/sherlock:v0.5.4 \
  dns run --config /app/config/config.yaml

# or

docker run --rm -it \
  ghcr.io/ch0ppy35/sherlock:v0.5.4 \
  dns test --server 1.1.1.1 \
       --host prom.example.com \
       --expected "10.0.0.1" \
       --type a
```
