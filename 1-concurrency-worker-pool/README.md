1. Concurrency in a distributed task queue / worker pool (very high probability)
They run security scans (distributed). You’ll likely get something like:

*“Write a Go worker pool that takes 100 scan jobs, processes them concurrently with max 10 workers, handles cancellations, and collects results. Then explain how you’d make it fault-tolerant across multiple pods in K8s.”*

What they test:

Goroutines + channels

context.Context for cancellation/timeout

sync.WaitGroup, errgroups, or worker pools

Graceful shutdown