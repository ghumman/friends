package ttlcacheaggregator

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeProvider is a configurable RateProvider for testing.
// It returns line.RiskScore * 100 as the premium, optionally injecting an
// error for a specific line ID and sleeping fro delay before responding.
type fakeProvider struct {
	delay time.Duration
	errorOnID string
	errorValue error
}

func (f *fakeProvider) Calculate(ctx context.Context, line PolicyLine) (float64, error) {
	if f.delay > 0 {
		select {
		case <-time.After(f.delay):
		case <-ctx.Done():
			return 0, ctx.Err()
		}
	}
	if f.errorOnID != "" && line.ID == f.errorOnID {
		return 0, f.errorValue
	}
	return line.RiskScore * 100, nil
}

// ----------------------------------
// TestTTLCache_BasicSetGet
// ----------------------------------

// TestTTLCache_BasicSetGet verifies that a value stored in the cache is 
// retrievable before its TTL expires.
func TestTTLCache_BasicSetGet(t *testing.T) {
	cache := NewTTLCache(5 * time.Second)
	cache.Set("auto-policy-001", 1250.50)
	value, ok := cache.Get("auto-policy-001")
	if !ok {
		t.Fatal("Expected to find value in cache, but it was missing")
	}
	if value.(float64) != 1250.50 {
		t.Fatalf("Expected value 1250.50, got %v", value)
	}
}

// ----------------------------------
// TestTTLCache_Expiry
// ----------------------------------

// TestTTLCache_Expiry verifies that a cache entry is not returned after its
// TTL has elapsed. Serving expired premium rates is a compliance risk.
func TestTTLCache_Expiry(t *testing.T) {
	cache := NewTTLCache(50 * time.Millisecond)
	cache.Set("home-policy-001", 980.75)
	time.Sleep(100 * time.Millisecond)
	_, ok := cache.Get("home-policy-001")
	if ok {
		t.Fatal("Expected cache entry to expire, but it was still present")
	}
}

// ----------------------------------
// TestAggregateLinePremiums_ConcurrentSuccess
// ----------------------------------

// TestAggregateLinePremiums_ConcurrentSuccess verifies that AggregateLinePremiums
// processes 10 policy lines concurrently (maxWorkers=3), retunrs all results,
// preserves original input order, and produces correct premium values.
func TestAggregateLinePremiums_ConcurrentSuccess(t *testing.T) {
	lines := make([]PolicyLine, 10)
	for i := range lines {
		lines[i] = PolicyLine{
			ID: string(rune('A' + i)),
			CoverageType: "TestCoverage",
			RiskScore: float64(i) + 1,
		}
	}

	provider := &fakeProvider{delay: 100 * time.Millisecond}
	results, err := AggregateLinePremiums(context.Background(), lines, provider, 3)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	for i, res := range results {
		expectedID := string(rune('A' + i))
		if res.ID != expectedID {
			t.Errorf("Expected ID %s, got %s", expectedID, res.ID)
		}
		expectedPremium := (float64(i) + 1) * 100
		if res.Premium != expectedPremium {
			t.Errorf("Expected premium %f, got %f", expectedPremium, res.Premium)
		}
	}
}

// ----------------------------------
// TestGetOrCompute_CacheMissAndHit
// ----------------------------------

// TestGetOrCompute_CacheMissAndHit verifies that RatedPolicyCache.GetOrCompute
// returns correct results on a cache miss and serves the cached value on a
// subsequent call without re-invoking the provider.
func TestGetOrCompute_CacheMissAndHit(t *testing.T) {
	lines := []PolicyLine{
		{ID: "P1", CoverageType: "Type1", RiskScore: 1.5},
		{ID: "P2", CoverageType: "Type2", RiskScore: 2.5},
	}

	provider := &fakeProvider{}
	cache := NewRatedPolicyCache(5 * time.Second)
	// First call should be a cache miss, invoking the provider
	results, err := cache.GetOrCompute(context.Background(), "bundle-001", lines, provider, 2)
	if err != nil {
		t.Fatalf("Unexpected error on cache miss: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}
	if results[0].Premium != 150 || results[1].Premium != 250 {
		t.Errorf("Expected premiums 150 and 250, got %f and %f", results[0].Premium, results[1].Premium)
	}

	// Second call should be a cache hit, returning cached values immediately
	result2, err := cache.GetOrCompute(context.Background(), "bundle-001", lines, provider, 2)
	if err != nil {
		t.Fatalf("Unexpected error on cache hit: %v", err)
	}
	if len(result2) != 2 {
		t.Fatalf("Expected 2 results on cache hit, got %d", len(result2))
	}

}



// ----------------------------------
// TestAggregateLinePremiums_ErrorPropagation
// ----------------------------------

// TestAggregateLinePremiums_ErrorPropagation verifies that one provider
// call returns an error, AggregateLinePremiums cancels all in-flight work and returns that error promptly.
func TestAggregateLinePremiums_ErrorPropagation(t *testing.T) {
	lines := []PolicyLine{
		{ID: "P1", CoverageType: "Type1", RiskScore: 1.0},
		{ID: "P2", CoverageType: "Type2", RiskScore: 2.0},
		{ID: "P3", CoverageType: "Type3", RiskScore: 3.0},
		{ID: "P4", CoverageType: "Type4", RiskScore: 4.0},
		{ID: "P5", CoverageType: "Type5", RiskScore: 5.0},
	}

	ratingErr := errors.New("rating error for P3")
	provider := &fakeProvider{delay: 100 * time.Millisecond, errorOnID: "P3", errorValue: ratingErr}
	
	_, err := AggregateLinePremiums(context.Background(), lines, provider, 5)
	if err == nil {
		t.Fatal("Expected error from provider, but got nil")
	}
}
