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

func downloadImage(imageLink string, downloadFolder string, splittedLink []string, wg *sync.WaitGroup) {
	defer wg.Done()

	imageFileName := splittedLink[len(splittedLink)-1]

	out, err := os.Create(fmt.Sprintf("%s/%s", strings.TrimRight(downloadFolder, "/"), imageFileName))
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer out.Close()

	resp, err := http.Get(imageLink)
	if err != nil {
		fmt.Println("Error downloading image:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Received non-200 status code %d for URL: %s\n", resp.StatusCode, imageLink)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println("Error saving image:", err)
		return
	}

	fmt.Printf("Downloaded %s\n", imageFileName)
}

func downloadImages(imagesList []string, downloadFolder string) {
	_, err := os.Stat(downloadFolder)
	if os.IsNotExist(err) {
		err := os.Mkdir(downloadFolder, 0755)
		if err != nil {
			fmt.Println("Error creating the download folder:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error checking if the download folder exists:", err)
		return
	}

	var wg sync.WaitGroup

	for _, imageLink := range imagesList {
		wg.Add(1)
		splittedLink := strings.Split(imageLink, "/")
		go downloadImage(imageLink, downloadFolder, splittedLink, &wg)
	}

	wg.Wait()
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
