# Health Check Sample

This sample demonstrates health checking and namespace operations against Temporal using mTLS authentication.

## Features

- Health check using `CheckHealth`
- Namespace description using `DescribeNamespace`
- mTLS client certificate authentication

## Running

```bash
go run main.go -client-cert /path/to/client.pem -client-key /path/to/client.key -namespace your-namespace
```

## Command Line Options

- `-target-host`: Host:port for the server (default: localhost:7233)
- `-namespace`: Namespace to use (default: default)
- `-client-cert`: Required path to client certificate
- `-client-key`: Required path to client private key
- `-server-root-ca-cert`: Optional path to root server CA cert
- `-server-name`: Server name for certificate verification
- `-insecure-skip-verify`: Skip certificate verification

## Testing Authentication

To compare the authentication between CheckHealth and DescribeNamespace try using a namespace you don't have access to:

```bash
go run main.go -client-cert /path/to/client.pem -client-key /path/to/client.key -namespace unauthorized-namespace
```

This will show permission denied errors when trying to describe the namespace but fine with check health.
