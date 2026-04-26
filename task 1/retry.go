package main

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"time"
)

const (
	maxRetries = 5
	baseDelay  = 500 * time.Millisecond
)

func IsRetryable(resp *http.Response, err error) bool {
	if err != nil {
		return true
	}

	if resp == nil {
		return false
	}

	switch resp.StatusCode {
	case 429, 500, 502, 503, 504:
		return true
	case 401, 404:
		return false
	default:
		return false
	}
}

func CalculateBackoff(attempt int) time.Duration {
	backoff := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
	jitter := time.Duration(rand.Int63n(int64(backoff)))
	return jitter
}

func ExecutePayment(ctx context.Context, client *http.Client, url string) error {
	for attempt := 0; attempt < maxRetries; attempt++ {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		resp, err := client.Do(req)

		if err == nil && resp.StatusCode == 200 {
			fmt.Println("✅ Success!")
			return nil
		}

		if !IsRetryable(resp, err) {
			return fmt.Errorf("❌ non-retryable error")
		}

		if attempt == maxRetries-1 {
			break
		}

		delay := CalculateBackoff(attempt)
		fmt.Printf("Attempt %d failed, waiting %v...\n", attempt+1, delay)

		select {
		case <-time.After(delay):
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return fmt.Errorf("❌ failed after retries")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	counter := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++

		if counter <= 3 {
			w.WriteHeader(503)
			fmt.Println("Server: 503")
			return
		}

		w.WriteHeader(200)
		fmt.Fprintln(w, `{"status":"success"}`)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := &http.Client{}

	err := ExecutePayment(ctx, client, server.URL)
	if err != nil {
		fmt.Println(err)
	}
}