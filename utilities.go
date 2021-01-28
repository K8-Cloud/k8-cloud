package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// referred from https://golangcode.com/download-a-file-from-a-url/

func DownloadFile(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func getFileFromURL(fileName string, fileUrl string)  {
	err := DownloadFile(fileName, fileUrl)
	if err != nil {
		panic(err)
	}
	fmt.Println("Downloaded: " + fileUrl)

}