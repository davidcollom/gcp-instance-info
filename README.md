# `gcp-instance-info`

GCP compute instance type and pricing information for Go, inspired by
[`LeanerCloud/ec2-instances-info`](https://github.com/LeanerCloud/ec2-instances-info).

The pricing and specification data for Google Cloud compute instances is embedded directly into your
binary, so you can look up vCPU count, memory, GPU availability and regional pricing without making
a network call or shipping a large data file alongside your application.

## Features

- Instance type metadata (vCPU, memory, family, GPU counts) embedded with `go:embed`
- Regional pricing for compute, storage and licenses
- A small updater to refresh the embedded data from the upstream source

## Requirements

- Go 1.22 or newer

## Installation

```bash
go get github.com/davidcollom/gcp-instance-info
```

## Usage

The entry point is `Data()`, which returns the parsed pricing data:

```go
import (
    gcpinstancesinfo "github.com/davidcollom/gcp-instance-info"
)

pricing, err := gcpinstancesinfo.Data()
if err != nil {
    return err
}

instance := pricing.Compute.Instances["n2-standard-4"]
fmt.Println(instance.VCPU, instance.Memory, instance.Family)
```

`Data()` returns a `GCPComputePricing` with four maps:

| Field | Description |
| ----- | ----------- |
| `Compute.Instances` | Instance type name -> `Instance` (family, vCPU, memory, GPU counts, per-region pricing). |
| `Compute.Storage` | Disk type -> `Storage` with monthly per-GB cost per region. |
| `Compute.Licenses` | License name (e.g. `rhel`, `windows`) -> `License` with per-region cost. |
| `Compute.Regions` | Region name -> `Region` (name and location). |

Each `Instance` exposes `Family`, `VCPU`, `Memory`, `GPU`, the per-GPU-type counts (`A100`,
`A10080GB`, `H10080GB`, `L4`) and a `Pricing` map of region -> `RegionPrices` (`Hour`, `HourSpot`,
`Month`, `Month1Y`, `Month2Y`, `Month3Y`, `MonthSpot`).

The embedded YAML lives in `data/instances.yaml` and is parsed on demand.

## Updating the data

The embedded data is generated from the upstream pricing file published by
[Cyclenerd/google-cloud-pricing-cost-calculator](https://github.com/Cyclenerd/google-cloud-pricing-cost-calculator).

Run the updater to refresh it:

```bash
go run ./cmd
```

This downloads the upstream `pricing.yml` from `SourceURL`, converts it, and writes the result to
`data/instances.yaml` (`DataPath`). The `//go:generate` directive on the embedded variable wires
this up automatically:

```bash
go generate ./...
```

The refresh interval used by the updater is 7 days (`defaultRefreshInterval`).

## Layout

| Path | Description |
| ---- | ----------- |
| `gcp_instance_data.go` | Types, the embedded dataset and the `Data()` accessor. |
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

No licence file is included in this repository.
