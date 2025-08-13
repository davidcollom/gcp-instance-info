package gcpinstancesinfo

import (
	_ "embed"
	"sync"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

//go:generate go run cmd/main.go
//go:embed data/instances.yaml
var staticDataBody []byte

var dataBody, backupDataBody []byte

var lock = sync.Mutex{}

//------------------------------------------------------------------------------

// Data generates the InstanceData object based on data sourced from
// gcpinstances.info. The data is available there as a JSON blob, which is
// converted into golang source-code and unmarshaled into a golang data
// structure by this library.
func Data() (out *GCPComputePricing, err error) {
	log.Println("Waiting for Mutex lock to read data...")
	lock.Lock()
	defer lock.Unlock()
	log.Printf("Dynamic data size: %d, static data size: %d, backup data size: %d",
		len(dataBody), len(staticDataBody), len(backupDataBody))

	if len(dataBody) == 0 {
		log.Printf("Dynamic data size is 0, using static data.")
		dataBody = staticDataBody
	}
	if len(dataBody) == 0 && len(staticDataBody) == 0 {
		log.Printf("Static data size is 0, using backup data.")
		dataBody = backupDataBody
	}
	if len(dataBody) == 0 {
		return nil, errors.New("no data available, static and backup data are both empty")
	}

	log.Printf("Data: [%d] %s", len(dataBody), dataBody)
	var in pricingYAML
	if err := yaml.Unmarshal(dataBody, &in); err != nil {
		log.Errorf("error unmarshaling data: %s", err)
		return nil, err
	}

	out = &GCPComputePricing{
		Compute: Compute{
			Instances: map[string]Instance{},
			Storage:   map[string]Storage{},
			Licenses:  map[string]License{},
			Regions:   map[string]Region{},
		},
	}
	log.Printf("Parsed %d instances, %d storage types, %d regions", len(in.Compute.Instances), len(in.Compute.Storage), len(in.Compute.Regions))

	// Regions
	for rk, r := range in.Compute.Regions {
		out.Compute.Regions[rk] = Region{
			Name:     firstNonEmpty(r.Name, rk),
			Location: r.Location,
		}
	}

	// Instances + per‑region prices
	for name, src := range in.Compute.Instances {
		inst := Instance{
			Family:       src.Family,
			InstanceType: name,
			VCPU:         src.CPU,
			Memory:       src.RAM,
			Pricing:      map[string]RegionPrices{},
		}

		// GPU summary if present
		for typ, n := range src.GPU {
			inst.GPU += n
			switch typ {
			case "a100":
				inst.A100 = &n
			case "a100-80gb":
				inst.A10080GB = &n
			case "h100-80gb":
				inst.H10080GB = &n
			case "l4":
				inst.L4 = &n
			}
		}

		for region, p := range src.Cost {
			inst.Pricing[region] = RegionPrices{
				Hour:      p.Hour,
				HourSpot:  p.HourSpot,
				Month:     p.Month,
				Month1Y:   p.Month1Y,
				Month3Y:   p.Month3Y,
				MonthSpot: p.MonthSpot,
			}
		}
		out.Compute.Instances[name] = inst

		// License add‑ons (windows/rhel/sles/…): src.Licenses[license][region] = prices
		for licName, byRegion := range src.Licenses {
			lic := out.Compute.Licenses[licName]
			if lic.Cost == nil {
				lic.Cost = map[string]RegionPrices{}
			}
			for region, lp := range byRegion {
				lic.Cost[region] = RegionPrices{
					Month:   lp.Month,
					Month1Y: lp.Month1Y,
					Month3Y: lp.Month3Y,
				}
			}
			out.Compute.Licenses[licName] = lic
		}
	}

	// Storage (if/when you use it)
	for sk, s := range in.Compute.Storage {
		st := Storage{
			Type: s.Type,
			Cost: map[string]RegionPrices{},
		}
		for region, p := range s.Cost {
			st.Cost[region] = RegionPrices{Month: p.Month}
		}
		out.Compute.Storage[sk] = st
	}

	return out, nil
}

func firstNonEmpty(ss ...string) string {
	for _, s := range ss {
		if s != "" {
			return s
		}
	}
	return ""
}
