package ttlcacheaggregator

import (
	"context"
	"sync"
	"time"
)

// PolicyLine represents a single line of coverage on an insurance policy.
type PolicyLine struct {
	ID string
	CoverageType string
	RiskScore float64
}

// LineResult holds the computed premium for a single plicy line. 
type LineResult struct {
	ID string
	Premium float64
}

// RateProvider is the interface for the external (expensive) premium
// calculation engine. Implementations may call remote rating services or 
// run actuarial models.
type RateProvider interface {
	Calculate (ctx context.Context, line PolicyLine) (float64, error)
}

// AggregateLinePremiums computes premiums for all lines concurrently, using at
// most maxWorkers goroutines. Results are returned in the same order as the
// input slice.
// 
// Requirements:
// 	- Concurrency is bounded by maxWorkers.
// 	- If any Calculate call returns an error, cancel all in-flight work and
// 	  return that error promptly
//	- Results must preserve the original input order.

func AggregateLinePremiums (ctx context.Context, lines []PolicyLine, provider RateProvider, maxWorkers int) ([]LineResult, error) {
	if maxWorkers <= 0 {
		maxWorkers = 1
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	results := make([]LineResult, len(lines))

	type job struct {
		index int
		line PolicyLine
	}

	jobs := make(chan job)

	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error

	worker := func() {
		defer wg.Done()
		for j := range jobs {
			select {
			case <- ctx.Done(): 
				return
			default:
			}

			result, err := provider.Calculate(ctx, j.line)
			if err != nil {
				once.Do( func() {
					firstErr = err
					cancel()
				})
				return
			}
			results[j.index] = LineResult{ID: j.line.ID, Premium: result}
		}
	}
	workerCount := maxWorkers
	if len(lines) < workerCount {
		workerCount = len(lines)
	}

	wg.Add(workerCount)
	for i:=0; i<workerCount; i++ {
		go worker()
	}

	// Feed jobs
	go func() {
		defer close(jobs)

		for i, line := range lines {
			select {
			case <- ctx.Done():
				return
			case jobs <- job{
				index: i,
				line: line,
			}:
			}
		}
	}()

	wg.Wait()

	if firstErr != nil {
		return nil, firstErr
	}

	if err := ctx.Err(); err != nil && err != context.Canceled {
		return  nil, err
	}
	return results, nil

}

// RatedPolicyCache wraps TTLCache with a concurrent-safe "compute once"
// layer. Multiple goroutines requesting the same policyKey at the same time
// must block until the first computation completes, then all receive the same
// result (singleflight behaviour)
// 
// TODO: add any fields you need to track in-flight computations and to 
// synchronise concurrent callers. Use stdlib sync primitives only.
type RatedPolicyCache struct {
	cache *TTLCache
	mu sync.Mutex
	// TODO: add inflight tracking fields here
}

// NewRatedPolicyCache returns a RatedPolicyCache backed by a TTLCache with
// the given time-to-live.
func NewRatedPolicyCache(ttl time.Duration) *RatedPolicyCache {
	return &RatedPolicyCache{
		cache: NewTTLCache(ttl),
	}
}

// GetOrCompute returns cached results for policyKey when available. On a cache
// miss it calls AggregateLinePremiums, stores the result, and returns it.
//
// Requirements:
// 	- Must not issue duplicate concurrent computations for the same policyKey
//	- The result is stored in the cache only when AggregateLinePremiums
//	  succeeds.

func (r *RatedPolicyCache) GetOrCompute(ctx context.Context, policyKey string, lines []PolicyLine, provider RateProvider, maxWorkers int) ([]LineResult, error){
	if cached, ok := r.cache.Get(policyKey); ok {
		return cached.([]LineResult), nil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if cached, ok := r.cache.Get(policyKey); ok {
		return cached.([]LineResult), nil
	}

	results, err := AggregateLinePremiums(ctx, lines, provider, maxWorkers)
	if err != nil {
		return nil, err
	}

	r.cache.Set(policyKey, results)
	return results, nil
}