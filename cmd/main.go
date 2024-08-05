package main

import (
	"fmt"
	"net/http"
	"os"

	gcpinstancesinfo "github.com/davidcollom/gcp-instance-info"
)

// Main function to download and update the embeded data
func main() {
	// Create the data directory if it does not exist
	if err := os.MkdirAll("data", os.ModePerm); err != nil {
		fmt.Printf("Failed to create data directory: %v\n", err)
		os.Exit(1)
	}

	updater := gcpinstancesinfo.NewDataFetcher(&http.Client{}, gcpinstancesinfo.SourceURL)

	// Download the file
	if err := updater.UpdateData(); err != nil {
		fmt.Printf("Failed Get new Data: %v\n", err)
		os.Exit(1)
	}

	// Create the file
	file, err := os.Create(gcpinstancesinfo.DataPath)
	if err != nil {
		fmt.Printf("Failed to create file: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	// Write the response body to the file
	_, err = file.Write(gcpinstancesinfo.GetData())
	if err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("File downloaded and saved to %s\n", gcpinstancesinfo.DataPath)

}
