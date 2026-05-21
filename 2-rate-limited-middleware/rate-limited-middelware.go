package ratelimitedmiddleware

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)


type RateLimitedRoundTripper struct {
	limiter *rate.Limiter
	next http.RoundTripper
	maxRetries int
}

func NewRateLimitedRoundTripper (rps float64, burst int, maxRetries int) *RateLimitedRoundTripper {
	return &RateLimitedRoundTripper{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
		next: http.DefaultTransport,
		maxRetries: maxRetries,
	}
}

func (r *RateLimitedRoundTripper) RoundTrip (req *http.Request) (*http.Response, error) {
	for attempt := 0; attempt<r.maxRetries; attempt++ {
		// wait for rate limit token
		err := r.limiter.Wait(req.Context())
		if err != nil {
			return nil, err
		}

		// make the request
		resp, err := r.next.RoundTrip(req)

		// check if successful
		if err == nil && resp.StatusCode < 400 {
			return resp, nil
		}

		// close body if it exists
		if resp != nil {
			resp.Body.Close()
		}

		// last attempt? Return error 
		if attempt == r.maxRetries {
			if err != nil {
				return nil, err
			}
			return nil, fmt.Errorf("failed with status %d", resp.StatusCode)
		}

		// calculate backoff
		delay := time.Duration(1 << attempt) * time.Second

		select {
		case <- time.After(delay):
		case <- req.Context().Done():
			return nil, req.Context().Err()
		}
	}

	return nil, fmt.Errorf("unreachable")
}

func main() {
	client := &http.Client {
		Transport: NewRateLimitedRoundTripper(5.0, 2, 3),
	}

	for i:=0; i<20; i++ {
		resp, err := client.Get("https://httpbin.org/get")
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else {
			fmt.Printf("Success! Status: %d\n", resp.StatusCode)
			resp.Body.Close()
		}
	}
}