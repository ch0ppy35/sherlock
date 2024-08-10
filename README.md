# sherlock

A simple tool to test your DNS records

## Usage

Have a config file configured with the tests you want to run

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

Run it

`sherlock run -c config/config.test.yaml`

Alternatively you can build and run it in Docker as a container

`docker run --rm -it -v "$(pwd)/config:/app/config" <IMAGENAME:TAG> --config /app/config/config.test.yaml`
