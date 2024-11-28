package gcpinstancesinfo

import (
	_ "embed"
	"strings"

	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v3"
)

type GCPComputePricing struct {
	Compute Compute `yaml:"compute"`
}

type Compute struct {
	Instances map[string]Instance `yaml:"instance"`
	Storage   map[string]Storage  `yaml:"storage"` // Map of String with name as key
	Licenses  map[string]License  // Map of String with Instance Type as key
	Regions   map[string]Region   // Map of String with Region as Key
}

type Region struct {
	Name     string
	Location string `yaml:"location"`
}

type Storage struct {
	Type string `yaml:"type"`
	// Costs are monthly per GB
	Cost map[string]RegionPrices `yaml:"cost"` // Map of String with Region as key
}

type License struct {
	Cost map[string]RegionPrices // Map of String with License (rhel,windows...) as key
}

type Instance struct {
	Family       string `yaml:"family"`
	InstanceType string `yaml:"type,omitempty"` // Essentially its name

	VCPU   float32 `yaml:"cpu"`
	Memory float32 `yaml:"ram"`
	GPU    int     // We don't get a direct GPU Count.. we need to calculate it

	// GPU Types
	A100     *int `yaml:"a100,omitempty"`
	A10080GB *int `yaml:"a100-80gb,omitempty"`
	H10080GB *int `yaml:"h100-80gb,omitempty"`
	L4       *int `yaml:"l4,omitempty"`

	// EnhancedNetworking bool   `json:"enhanced_networking"`
	// ECURaw             json.RawMessage `json:"ECU"`
	// ECU                string
	// VCPU               int
	// PhysicalProcessor  string                  `json:"physical_processor"`
	// Generation         string                  `json:"generation"`
	// EBSIOPS            float32                 `json:"ebs_iops"`
	// NetworkPerformance string                  `json:"network_performance"`
	// EBSThroughput      float32                 `json:"ebs_throughput"`
	// PrettyName         string                  `json:"pretty_name"`

	Pricing map[string]RegionPrices `yaml:"cost"`
}

type RegionPrices struct {
	Hour      float64 `yaml:"hour"`
	HourSpot  float64 `yaml:"hour_spot,omitempty"`
	Month     float64 `yaml:"month,omitempty"`
	Month1Y   float64 `yaml:"month_1y,omitempty"`
	Month2Y   float64 `yaml:"month_2y,omitempty"`
	Month3Y   float64 `yaml:"month_3y,omitempty"`
	MonthSpot float64 `yaml:"month_spot,omitempty"`
}

//go:generate go run cmd/main.go
//go:embed data/instances.yaml
var staticDataBody []byte

var dataBody, backupDataBody []byte

//------------------------------------------------------------------------------

// Data generates the InstanceData object based on data sourced from
// gcpinstances.info. The data is available there as a JSON blob, which is
// converted into golang source-code and unmarshaled into a golang data
// structure by this library.
func Data() (*GCPComputePricing, error) {

	var j GCPComputePricing

	if len(dataBody) > 0 {
		log.Println("We have updated data, trying to unmarshal it")
		err := yaml.Unmarshal(dataBody, &j)
		if err != nil {
			log.Printf("couldn't unmarshal the updated data, reverting to the backup data : %s", err.Error())
			err := yaml.Unmarshal(backupDataBody, &j)
			if err != nil {
				return nil, errors.Errorf("couldn't unmarshal backup data: %s", err.Error())
			}
			backupDataBody = []byte{}
		}
	} else {
		log.Println("Using the static instance type data")
		err := yaml.Unmarshal(staticDataBody, &j)
		if err != nil {
			return nil, errors.Errorf("couldn't unmarshal data: %s", err.Error())
		}
	}

	// We Need to populate InstanceType for all instances
	for f, i := range j.Compute.Instances {
		if i.InstanceType == "" {
			i.InstanceType = f
		}
		// Set the Family from the first portion of the InstanceType
		i.Family = strings.Split(i.InstanceType, "-")[0]

		// We also need to find the GPU's
		if i.A100 != nil {
			i.GPU = *i.A100
		}
		if i.A10080GB != nil {
			i.GPU = *i.A10080GB
		}
		if i.H10080GB != nil {
			i.GPU = *i.H10080GB
		}
		if i.L4 != nil {
			i.GPU = *i.L4
		}
		j.Compute.Instances[f] = i
	}

	// for family, instanceTypes := range j {
	// 	// TODO: Find better way of detecting "Licenses"
	// 	if reflect.TypeOf(family).String() != "jsonInstance" {
	// 		log.Println("Ignoring %s as type is %T", family, family)
	// 		continue
	// 	}
	// 	for instance := range instanceTypes {

	// 	}
	// }

	// sort.Slice(j, func(a, b string) bool {
	// 	// extract the instance family, such as "c5" for "c5.large"
	// 	family_i := strings.Split(j.Compute.Instances[a].InstanceType, ".")[0]
	// 	family_j := strings.Split(j.Compute.Instances[b].InstanceType, ".")[0]

	// 	// we first compare only based on the family
	// 	switch strings.Compare(family_i, family_j) {
	// 	case -1:
	// 		return true
	// 	case 1:
	// 		return false
	// 	}

	// 	// within the same family we compare by memory size, but always keeping metal instances last
	// 	return d[a].Memory < d[b].Memory || strings.HasSuffix(j[a].InstanceType, "metal")
	// })

	return &j, nil
}
