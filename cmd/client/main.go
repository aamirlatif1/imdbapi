package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type result struct {
	id       int
	status   int
	body     string
	duration time.Duration
	err      error
}

func main() {
	url := flag.String("url", "http://localhost:4000/v1/movies/3", "target URL")
	n := flag.Int("n", 10, "number of concurrent requests")
	payload := flag.String("data", `{"runtime":100,"Title":"Test"}`, "JSON request body")
	//delay := flag.Duration("delay", 1*time.Millisecond, "delay between each request")
	flag.Parse()

	client := &http.Client{Timeout: 10 * time.Second}

	var wg sync.WaitGroup
	start := make(chan struct{})
	results := make(chan result, *n)

	for i := 1; i <= *n; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			<-start
			// stagger requests so request i is sent (i-1)*delay after the first
			//time.Sleep(time.Duration(id-1) * *delay)
			results <- send(client, id, *url, *payload)
		}(i)
	}

	close(start)
	wg.Wait()
	close(results)

	counts := make(map[int]int)
	for r := range results {
		if r.err != nil {
			log.Printf("request %2d: error: %v", r.id, r.err)
			counts[0]++
			continue
		}
		fmt.Printf("request %2d: %d (%v) %s\n", r.id, r.status, r.duration, r.body)
		counts[r.status]++
	}

	fmt.Println("\nsummary:")
	for status, c := range counts {
		if status == 0 {
			fmt.Printf("  errors: %d\n", c)
			continue
		}
		fmt.Printf("  %d %s: %d\n", status, http.StatusText(status), c)
	}
}

func send(client *http.Client, id int, url, payload string) result {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result{id: id, err: err}
	}
	req.Header.Set("Content-Type", "application/json")

	begin := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return result{id: id, err: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return result{id: id, err: err}
	}
	return result{
		id:       id,
		status:   resp.StatusCode,
		body:     string(bytes.TrimSpace(body)),
		duration: time.Since(begin),
	}
}
