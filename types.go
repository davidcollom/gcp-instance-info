package gcpinstancesinfo

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

// Input model — mirrors pricing.yml structure (subset)
type pricingYAML struct {
	Compute struct {
		Instances map[string]struct {
			Family string  `yaml:"family"`
			CPU    float32 `yaml:"cpu"` // vCPUs
			RAM    float32 `yaml:"ram"` // GB

			// Per‑region prices for this instance type
			Cost map[string]struct {
				Hour      float64 `yaml:"hour"`
				HourSpot  float64 `yaml:"hour_spot"`
				Month     float64 `yaml:"month"`
				Month1Y   float64 `yaml:"month_1y"`
				Month3Y   float64 `yaml:"month_3y"`
				MonthSpot float64 `yaml:"month_spot"`
			} `yaml:"cost"`

			// Optional license add‑ons per region (windows, rhel, sles, sles-sap…)
			Licenses map[string]map[string]struct {
				Month   float64 `yaml:"month"`
				Month1Y float64 `yaml:"month_1y"`
				Month3Y float64 `yaml:"month_3y"`
			} `yaml:"licenses,omitempty"`

			// Optional accelerator summary, if present
			GPU map[string]int `yaml:"gpu,omitempty"` // e.g. {"a100": 1, "l4": 2}
		} `yaml:"instance,omitempty"`

		// Regions index
		Regions map[string]struct {
			Location string `yaml:"location"`
			Name     string `yaml:"name,omitempty"`
		} `yaml:"regions,omitempty"`

		// Storage price table (simplified)
		Storage map[string]struct {
			Type string              `yaml:"type"`
			Cost map[string]struct { // by region
				Month float64 `yaml:"month"`
			} `yaml:"cost"`
		} `yaml:"storage,omitempty"`
	} `yaml:"compute"`
}
