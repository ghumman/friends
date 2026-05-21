2. Rate-limited API client or middleware (high probability)
Given security testing involves external APIs (cloud providers, customer APIs).

*“Implement a rate-limited HTTP client in Go that allows 5 requests per second, with burst of 2, and retries with backoff on 429/503.”*

What they test:

golang.org/x/time/rate

http.RoundTripper customization

Exponential backoff (not just a loop sleep)

