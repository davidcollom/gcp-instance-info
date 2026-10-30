# `gcp-instance-info`

GCP instance type information for Go, inspired by [`LeanerCloud/ec2-instances-info`](https://github.com/LeanerCloud/ec2-instances-info).

The pricing and specification data for Google Cloud compute instances is embedded directly into your
binary, so you can look up vCPU count, memory, GPU availability and regional pricing without making a
network call or shipping a large data file alongside your application.

## Features

- Instance type metadata (vCPU, memory, family, GPU counts) embedded with `go:embed`
- Regional pricing for compute, storage and licenses
- A small updater to refresh the embedded data from the upstream source

## Requirements

- Go 1.21 or newer

## Installation

```bash
go get github.com/davidcollom/gcp-instance-info
```

## Usage

```go
import (
    gcpinstancesinfo "github.com/davidcollom/gcp-instance-info"
)

instance, err := gcpinstancesinfo.GetInstance("n2-standard-4")
if err != nil {
    return err
}

fmt.Println(instance.VCPU, instance.Memory, instance.Family)
```

The embedded YAML is parsed on demand from `data/instances.yaml`.

## Updating the data

The embedded data is generated from the upstream pricing file published by
[Cyclenerd/google-cloud-pricing-cost-calculator](https://github.com/Cyclenerd/google-cloud-pricing-cost-calculator).

Run the updater to refresh it:

```bash
go run ./cmd
```

This downloads the upstream `pricing.yml`, extracts the instance, storage and license pricing, and
writes the result to `data/instances.yaml`.

The refresh interval used by the updater is 7 days.

## Layout

| Path | Description |
| ---- | ----------- |
| `gcp_instance_data.go` | Types and embedded data accessors. |
| `updater.go` | Downloads and transforms the upstream pricing data. |
| `logger.go` | Logging helper used by the updater. |
| `cmd/main.go` | Command that refreshes `data/instances.yaml`. |
| `data/instances.yaml` | The generated, embedded pricing and specification data. |
| `.github/workflows/go.yaml` | CI workflow running the Go tests. |

## Development

```bash
# Run the tests
go test ./...

# Lint
golangci-lint run
```

## License

See the repository for licensing details.
