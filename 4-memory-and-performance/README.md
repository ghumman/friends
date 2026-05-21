4. Memory & performance debugging (medium, but senior-level)
Synack cares about efficiency (long-running security scans).

“You see memory growing over time in a Go service. How do you debug it? Show me pprof in action conceptually. What’s a common cause of unbounded growth in Go you’ve seen?”

Possible follow-up:

“We have a struct with a []byte slice. After slicing it to a smaller length, why doesn’t memory shrink? How do you force it?”