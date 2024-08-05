package gcpinstancesinfo

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	SourceURL              = "https://raw.githubusercontent.com/Cyclenerd/google-cloud-pricing-cost-calculator/master/pricing.yml"
	DataPath               = "./data/instances.yaml"
	defaultRefreshInterval = 7
)

var httpClient = &http.Client{}

type DataFetcher struct {
	client    *http.Client
	sourceURL string
}

func NewDataFetcher(client *http.Client, sourceURL string) *DataFetcher {
	return &DataFetcher{
		client:    client,
		sourceURL: sourceURL,
	}
}

func (df *DataFetcher) fetchData() ([]byte, error) {
	req, err := http.NewRequest("GET", df.sourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	resp, err := df.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %s", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}
	return body, nil
}

func (df *DataFetcher) UpdateData() error {
	log.Printf("Dynamic data size before: %d, downloading new instance type data.", len(dataBody))

	body, err := df.fetchData()
	if err != nil {
		return err
	}

	// Save our previous data (this may be the static version)
	if len(dataBody) > 0 {
		backupDataBody = dataBody
	} else {
		backupDataBody = staticDataBody
	}

	dataBody = body
	log.Println("Data size after:", len(dataBody))

	return nil
}

func GetData() []byte {
	return dataBody
}

func (df *DataFetcher) Updater(refreshDays int) error {
	if refreshDays <= 0 {
		refreshDays = defaultRefreshInterval
	}
	refreshInterval := time.Duration(refreshDays) * 24 * time.Hour

	if err := df.UpdateData(); err != nil {
		log.Printf("Failed to download updated data: %s", err.Error())
		return fmt.Errorf("error downloading new data: %s", err)
	}

	// refresh the data every refreshInterval
	ticker := time.NewTicker(refreshInterval)
	defer ticker.Stop()

	for range ticker.C {
		err := df.UpdateData()
		if err != nil {
			log.Println("Error refreshing data:", err)
		}
	}
	return nil
}
