package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

func downloadImage(imageLink string, downloadFolder string, splittedLink []string, wg *sync.WaitGroup) error {
	defer wg.Done()

	imageFileName := splittedLink[len(splittedLink)-1]

	out, err := os.Create(fmt.Sprintf("%s/%s", strings.TrimRight(downloadFolder, "/"), imageFileName))
	if err != nil {
		return fmt.Errorf("Error creating file:", err)
	}
	defer out.Close()

	resp, err := http.Get(imageLink)
	if err != nil {
		return fmt.Errorf("Error downloading image:", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error: Received non-200 status code %d for URL: %s\n", resp.StatusCode, imageLink)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("Error saving image:", err)
	}

	fmt.Printf("Downloaded %s\n", imageFileName)

	return nil
}

func downloadImages(imagesList []string, downloadFolder string) error {
	_, err := os.Stat(downloadFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(downloadFolder, 0755)
		if err != nil {
			return fmt.Errorf("Error creating the download folder:", err)
		}
	} else if err != nil {
		return fmt.Errorf("Error checking if the download folder exists:", err)
	}

	var wg sync.WaitGroup

	for _, imageLink := range imagesList {
		wg.Add(1)
		splittedLink := strings.Split(imageLink, "/")
		go downloadImage(imageLink, downloadFolder, splittedLink, &wg)
	}

	wg.Wait()

	return nil
}

func main() {
	imagesLinksFile, downloadFolder := "images.txt", "downloadedImages"

	file, err := os.Open(imagesLinksFile)
	if err != nil {
		fmt.Printf("Error opening the images list file: %s \n", err)
		return
	}
	defer file.Close()

	var imagesList []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		imagesList = append(imagesList, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading the images list file: %s \n", err)
		return
	}

	downloadImages(imagesList, downloadFolder)
}
