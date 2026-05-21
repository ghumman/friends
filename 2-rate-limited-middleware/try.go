package ratelimitedmiddleware

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/time/rate"
)

type RateLimiter struct {
    limiter   *rate.Limiter
    next      http.RoundTripper
    maxTries  int
}

func NewRateLimiter(rps float64, burst int, maxTries int) *RateLimiter {
    return &RateLimiter{
        limiter:  rate.NewLimiter(rate.Limit(rps), burst),
        next:     http.DefaultTransport,
        maxTries: maxTries,
    }
}

func (r *RateLimiter) RoundTrip(req *http.Request) (*http.Response, error) {
    for attempt := 0; attempt <= r.maxTries; attempt++ {
        // Rate limit before each attempt
        if err := r.limiter.Wait(req.Context()); err != nil {
            return nil, err
        }
        
        resp, err := r.next.RoundTrip(req)
        
        // Handle successful response (no network error)
        if err == nil {
            // Success (2xx)
            if resp.StatusCode >= 200 && resp.StatusCode < 300 {
                return resp, nil
            }
            
            // ✅ Fix 1: Only retry on 429 or 503
            if resp.StatusCode == 429 || resp.StatusCode == 503 {
                resp.Body.Close()  // Close before retry
                // Continue to retry logic
            } else {
                // Non-retryable status - return immediately
                return resp, nil
            }
        } else {
            // ✅ Fix 2: Network error - resp is nil, nothing to close
            // Just continue to retry
        }
        
        // ✅ Fix 3: Don't sleep on last attempt
        if attempt == r.maxTries {
            break
        }
        
        // Exponential backoff: 1s, 2s, 4s, 8s...
        delay := time.Duration(1<<attempt) * time.Second
        time.Sleep(delay)
    }
    
    return nil, fmt.Errorf("not able to reach the server after %d attempts", r.maxTries+1)
}

func main() {
    client := &http.Client{
        Transport: NewRateLimiter(5, 2, 3),
    }
    
    for i := 0; i < 20; i++ {
        resp, err := client.Get("https://httpbin.org/get")
        if err != nil {
            fmt.Printf("Request %2d: ❌ FAILURE: %v\n", i+1, err)
        } else {
            fmt.Printf("Request %2d: ✅ SUCCESS (status: %d)\n", i+1, resp.StatusCode)
            resp.Body.Close()
        }
    }
}