# Sherlock

A CLI tool for testing DNS records.

## Overview

Sherlock is a command-line tool designed to perform DNS record tests based on a specified configuration file. It allows you to run various types of DNS checks, such as verifying A, AAAA, CNAME, MX, TXT, and NS records.

This tool supports a configuration file in YAML format where you can define the expected DNS records for different hosts. The `run` command executes all tests defined in the configuration and provides a summary of any discrepancies found.

## Configuration

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

## Usage

### Running Tests

Execute the DNS tests defined in your configuration file using the following command:

```bash
sherlock run --config path/to/config.yaml
```

Replace `path/to/config.yaml` with the actual path to your configuration file.

### Docker

Alternatively, you can run Sherlock inside a Docker container. First, ensure you have a Docker image built or available. Then, execute:

```bash
docker run --rm -it -v "$(pwd)/config:/app/config" ghcr.io/ch0ppy35/sherlock:v0.2.0 run --config /app/config/config.yaml
```
