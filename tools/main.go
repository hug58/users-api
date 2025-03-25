package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

func makeRequest(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error making request to %s: %s\n", url, err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request to %s completed with status: %s\n", url, resp.Status)
}

func main() {
	url := "http://127.0.0.1:40955/api/v1/users" // Cambia esta URL por la que desees
	requestsPerSecond := 10000
	var wg sync.WaitGroup

	for i := 0; i < requestsPerSecond; i++ {
		wg.Add(1)
		go makeRequest(url, &wg)
		time.Sleep(time.Second / time.Duration(requestsPerSecond))
	}

	wg.Wait()
	fmt.Println("All requests completed.")
}
