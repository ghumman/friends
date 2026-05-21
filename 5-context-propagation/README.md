5. Context propagation in distributed systems (very high for “distributed systems” requirement)
“We have a gRPC call from Service A → B → C. How do you pass a deadline/trace ID from A to C in Go? Write pseudocode. What happens if B forgets to pass the context?”

They’ll look for:

context.WithTimeout / WithDeadline

Metadata propagation (gRPC, HTTP headers)

Cancellation cascading