package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

func main() {
	urls := []string{
		"https://www.google.com/",
		"https://www.google.com/",
		"https://www.google.com/",
		"https://www.google.com/",
	}

	success := sendRequests(urls)
	if success {
		fmt.Println("all of the requests have received successful responses")
	} else {
		fmt.Println("some of the requests have unsuccessful responses")
	}
}

func sendRequests(urls []string) bool {
	successCount := uint32(0)
	wg := sync.WaitGroup{}

	for _, url := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			resp, err := sendRequestWithTimeout(url)
			if err != nil {
				log.Println("send request err:", err.Error())
				return
			}
			if resp.StatusCode == http.StatusOK {
				atomic.AddUint32(&successCount, 1)
			}
		}(url)
	}

	wg.Wait()
	return uint32(len(urls)) == successCount
}

func sendRequestWithTimeout(url string) (*http.Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	return http.DefaultClient.Do(req)
}
