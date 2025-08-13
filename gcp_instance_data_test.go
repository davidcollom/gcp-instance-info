package gcpinstancesinfo

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testRegion = "us-east1"

var testMu = sync.Mutex{}

// Mock data for testing
var mockStaticData = `
compute:
  instance:
    n1-standard-1:
      cpu: 1
      ram: 3.75
      cost:
        us-east1:
          hour: 0.0475
`

var mockUpdatedData = `
compute:
  instance:
    n1-standard-2:
      cpu: 2
      ram: 7.5
      cost:
        us-east1:
          hour: 0.095
`

var mockBackupData = `
compute:
  instance:
    n1-standard-4:
      cpu: 4
      ram: 15
      cost:
        us-east1:
          hour: 0.19
`

func TestData(t *testing.T) {
	tests := []struct {
		name             string
		dataBody         []byte
		backupDataBody   []byte
		staticDataBody   []byte
		expectedInstance string
		expectedVCPU     float32
		expectedMemory   float32
		expectedHour     float64
		expectError      bool
	}{
		{
			name:             "StaticData",
			dataBody:         nil,
			backupDataBody:   nil,
			staticDataBody:   []byte(mockStaticData),
			expectedInstance: "n1-standard-1",
			expectedVCPU:     1,
			expectedMemory:   3.75,
			expectedHour:     0.0475,
			expectError:      false,
		},
		{
			name:             "UpdatedData",
			dataBody:         []byte(mockUpdatedData),
			backupDataBody:   nil,
			staticDataBody:   []byte(mockStaticData),
			expectedInstance: "n1-standard-2",
			expectedVCPU:     2,
			expectedMemory:   7.5,
			expectedHour:     0.095,
			expectError:      false,
		},
		{
			name:             "BackupData",
			dataBody:         []byte{},
			backupDataBody:   []byte(mockBackupData),
			staticDataBody:   []byte{},
			expectedInstance: "n1-standard-4",
			expectedVCPU:     4,
			expectedMemory:   15,
			expectedHour:     0.19,
			expectError:      false,
		},
		{
			name:             "BackupDataFail",
			dataBody:         []byte{},
			backupDataBody:   []byte("invalid backup data"),
			staticDataBody:   []byte{},
			expectedInstance: "",
			expectedVCPU:     0,
			expectedMemory:   0,
			expectedHour:     0,
			expectError:      true,
		},
		{
			name:           "FamilyAndGPU",
			dataBody:       nil,
			backupDataBody: nil,
			staticDataBody: []byte(`
compute:
  instance:
    n1-standard-1:
      type: n1-standard-1
      a100: 1
      cost:
        us-east1:
          hour: 0.0475
`),
			expectedInstance: "n1-standard-1",
			expectedVCPU:     0,
			expectedMemory:   0,
			expectedHour:     0.0475,
			expectError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()
			testMu.Lock()

			dataBody = tt.dataBody
			backupDataBody = tt.backupDataBody
			staticDataBody = tt.staticDataBody

			pricing, err := Data()
			testMu.Unlock()

			if tt.expectError {
				require.Error(t, err)
				assert.Nil(t, pricing)
			} else {
				require.NoError(t, err)
				require.NotNil(t, pricing)
				// require.NotNil(t, pricing.Compute)
				require.GreaterOrEqual(t, len(pricing.Compute.Instances), 1, "Expected at least one instance in pricing data: %+v", pricing)

				instance, exists := pricing.Compute.Instances[tt.expectedInstance]
				require.True(t, exists, "Instance %s should exist: %+v", tt.expectedInstance, pricing.Compute.Instances)
				assert.Equal(t, tt.expectedVCPU, instance.VCPU)
				assert.Equal(t, tt.expectedMemory, instance.Memory)
				assert.Equal(t, tt.expectedHour, instance.Pricing[testRegion].Hour)
			}
		})
	}
}
