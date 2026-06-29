package main

import (
	"io"
	"net/http"
	"sync"
)

type LoadBalancer interface {
	GetNextBackend() string
	ProxyHandler(w http.ResponseWriter, r *http.Request) 
}

type RoundRobinProxy struct {
	backends []string
	backendLength int
	currentIndex int
	mu sync.Mutex
}

func NewRoundRobinProxy(backends []string) *RoundRobinProxy {
	return & RoundRobinProxy{
		backends: backends,
		backendLength: len(backends),
		currentIndex: 0,
	}
}

func (p *RoundRobinProxy) GetNextBackend()string {
	p.mu.Lock()
	defer p.mu.Unlock()

	backend := p.backends[p.currentIndex]
	p.currentIndex = (p.currentIndex + 1) % p.backendLength
	return backend
}

func (p *RoundRobinProxy) ProxyHandler(w http.ResponseWriter, r *http.Request) {
	backend := p.GetNextBackend()

	proxyReq, err := http.NewRequest(r.Method, backend + r.URL.Path, r.Body)
	if err != nil {
		http.Error(w, "Error creating proxy request", http.StatusInternalServerError)
		return
	}

	for key, values := range r.Header {
		for _, val := range values {
			proxyReq.Header.Add(key, val)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		http.Error(w, "Backend call failed", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, val := range values {
			w.Header().Add(key, val)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	targets := []string{
		"http://10.0.0.1:8080",
		"http://10.0.0.2:8080",
		"http://10.0.0.3:8080",
	}
	_ = NewRoundRobinProxy(targets)
}