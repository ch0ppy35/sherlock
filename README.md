# Sherlock

A CLI tool for testing DNS records.

## Overview

Sherlock is a command-line tool designed to perform DNS record tests based on a specified configuration file or individual params. It allows you to run various types of DNS checks, such as verifying A, AAAA, CNAME, MX, TXT, and NS records.

This tool supports a configuration file in YAML format where you can define the expected DNS records for different hosts. The 'run' command executes all tests defined in the configuration and provides a summary of any discrepancies found. Alternatively, you can use the 'test' command to query a DNS server for a specific record type and compare the results with expected values directly from the command line.

This tool is generally intended to be used within a container, making it ideal for integration into CI/CD pipelines or scheduled tasks like cron jobs running in Kubernetes. Binaries are also provided in the GitHub release, or you can build the binary locally using the `make` command.

See `sherlock -h` for more info

## Usage

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
sherlock run --config path/to/config.yaml

# or

sherlock test --server 1.1.1.1 --host prom.example.com --expected "10.0.0.1" --type a
```

### Docker

Alternatively, you can run Sherlock inside a Docker container

```bash
docker run --rm -it \
  -v "$(pwd)/config:/app/config" \
  ghcr.io/ch0ppy35/sherlock:v0.5.0 \
  run --config /app/config/config.yaml

# or

docker run --rm -it \
  ghcr.io/ch0ppy35/sherlock:v0.5.0 \
  test --server 1.1.1.1 \
       --host prom.example.com \
       --expected "10.0.0.1" \
       --type a

```
