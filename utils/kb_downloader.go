package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func DownloadIfChanged(localPath, localETagPath, remoteURL string) error {
	// Check if the local file exists
	file, err := os.Open(localPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	defer file.Close()

	// Read the existing ETag from the separate metadata file
	localETag, err := readETag(localETagPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Make a GET request to the remote URL with the existing ETag
	req, err := http.NewRequest("GET", remoteURL, nil)
	if err != nil {
		return err
	}

	if localETag != "" {
		req.Header.Add("If-None-Match", localETag)
	}

	// Make the GET request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check if the response status code indicates success
	if resp.StatusCode == http.StatusNotModified {
		// fmt.Println("File is up to date. No need to download.")
		return nil
	} else if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %s", resp.Status)
	}

	// Create or open the local file for writing
	file, err = os.Create(localPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Copy the response body to the local file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	// Save the received ETag to the separate metadata file
	err = saveETag(localETagPath, resp.Header.Get("ETag"))
	if err != nil {
		return err
	}

	fmt.Println("File downloaded successfully.")
	return nil
}

func readETag(etagPath string) (string, error) {
	// Read the ETag from the separate metadata file
	file, err := os.Open(etagPath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var etag string
	_, err = fmt.Fscanln(file, &etag)
	if err != nil && err != io.EOF {
		return "", err
	}

	return etag, nil
}

func saveETag(etagPath, etag string) error {
	// Save the ETag to the separate metadata file
	file, err := os.Create(etagPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the ETag to the file
	_, err = file.WriteString(etag)
	return err
}
