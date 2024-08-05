package gcpinstancesinfo

import (
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockDataFetcher struct {
	client        *http.Client
	sourceURL     string
	fetchDataFunc func() ([]byte, error)
}

func (m *MockDataFetcher) fetchData() ([]byte, error) {
	return m.fetchDataFunc()
}

func (m *MockDataFetcher) UpdateData() error {
	body, err := m.fetchData()
	if err != nil {
		return err
	}

	if len(dataBody) > 0 {
		backupDataBody = dataBody
	} else {
		backupDataBody = staticDataBody
	}

	dataBody = body
	return nil
}

func (m *MockDataFetcher) Updater(refreshDays int) error {
	if refreshDays <= 0 {
		refreshDays = defaultRefreshInterval
	}
	refreshInterval := time.Duration(refreshDays) * 24 * time.Hour

	if err := m.UpdateData(); err != nil {
		return err
	}

	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := m.UpdateData()
		if err != nil {
			return err
		}
	}
	return nil
}

func TestUpdateData(t *testing.T) {
	tests := []struct {
		name           string
		serverResponse string
		serverStatus   int
		expectedError  bool
	}{
		{
			name:           "SuccessfulDownload",
			serverResponse: "mock data",
			serverStatus:   http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "ErrorSendingRequest",
			serverResponse: "",
			serverStatus:   http.StatusInternalServerError,
			expectedError:  true,
		},
		{
			name:           "ErrorReadingResponseBody",
			serverResponse: "",
			serverStatus:   http.StatusOK,
			expectedError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.serverStatus != http.StatusOK {
					http.Error(w, "error", tt.serverStatus)
					return
				}
				w.WriteHeader(tt.serverStatus)
				_, err := io.WriteString(w, tt.serverResponse)
				require.NoError(t, err)
			}))
			defer server.Close()

			mockFetcher := &MockDataFetcher{
				client:    httpClient,
				sourceURL: server.URL,
				fetchDataFunc: func() ([]byte, error) {
					resp, err := http.Get(server.URL)
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()
					if tt.serverStatus != http.StatusOK {
						return nil, errors.New("error fetching data")
					}
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						return nil, err
					}
					return body, nil
				},
			}

			dataBody = nil
			backupDataBody = nil

			err := mockFetcher.UpdateData()

			if tt.expectedError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, []byte(tt.serverResponse), dataBody)
			}
		})
	}
}
