package main

import (
	"sync"
	"time"
)


type TokenBucketLimiter struct {
	maxTokens float64
	currentTokens float64
	refillRate float64
	lastTimestamp time.Time
	mu sync.Mutex
}


func NewTokenBucket( maxTokens float64, refillRate float64) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		maxTokens: maxTokens,
		currentTokens: maxTokens,
		refillRate: refillRate,
		lastTimestamp: time.Now(),
	}
}


func (tb *TokenBucketLimiter) Allow()bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	currentTime := time.Now()
	elapsedTime := currentTime.Sub(tb.lastTimestamp).Seconds()

	newTokens := elapsedTime * tb.refillRate

	tb.currentTokens += newTokens

	if tb.currentTokens > tb.maxTokens {
		tb.currentTokens = tb.maxTokens
	}

	if tb.currentTokens >= 1.0 {
		tb.currentTokens--
		return true
	}
	return false
	
}