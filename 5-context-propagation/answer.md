Here is your ultra-focused, high-impact interview cheat sheet for context propagation.

---

## 1. The Core Concepts (The Talk)

When the interviewer asks how to pass a deadline or trace ID from Service A to C, they want to hear three specific concepts:

* **Metadata Injection/Extraction:** You don't pass `context.Context` directly over the network wire. Instead, gRPC client interceptors **inject** the context data (deadlines, trace IDs) into gRPC **metadata** (which maps to HTTP/2 headers). Service B's server interceptor then **extracts** this metadata back into a fresh Go `context.Context`.
* **Cancellation Cascading:** If Service A times out or cancels its context, the underlying gRPC connection drops. Service B detects this on `ctx.Done()`, stops working, and cancels its outgoing call to Service C. This prevents wasted resources (zombie requests).
* **Automatic gRPC Deadline Propagation:** By default, if Service B uses the incoming `ctx` from Service A to make the call to Service C, **gRPC automatically adjusts and passes the remaining deadline time** to Service C.

---

## 2. The Code (The Walk)

Here is the clean, idiomatic Go pseudocode demonstrating how Service B must correctly accept the incoming context and forward it to Service C.

```go
// Service B implementation
func (s *ServiceB) DoSomething(ctx context.Context, req *pb.Request) (*pb.Response, error) {
    // 1. (Optional) Explicitly add/modify gRPC Metadata if injecting custom Trace IDs
    // md := metadata.Pairs("x-trace-id", "12345")
    // ctx = metadata.NewOutgoingContext(ctx, md)

    // 2. Pass the incoming 'ctx' directly to the next gRPC call. 
    // This automatically forwards deadlines and cascaded cancellations.
    resC, err := s.ClientC.CallServiceC(ctx, &pb.RequestC{Data: req.Data})
    if err != nil {
        return nil, err
    }

    return &pb.Response{Data: resC.Data}, nil
}

```

If Service A wants to set a specific timeout right at the start:

```go
// In Service A
ctx, cancel := context.WithTimeout(context.Background(), 2 * time.Second)
defer cancel() // Always clean up resources

resB, err := clientB.DoSomething(ctx, &req)

```

---

## 3. What happens if B forgets to pass the context?

If Service B mistakenly passes `context.Background()` or `context.TODO()` to Service C instead of the incoming `ctx`, **the distributed chain breaks completely.**

Here is what goes wrong:

| Problem | Consequence |
| --- | --- |
| **Trace ID Broken** | Service C receives no trace context. In your logs/APM dashboard (like Jaeger or OpenTelemetry), the request path looks like two unrelated fragments instead of one unified trace. |
| **No Cancellation Cascading** | If Service A times out and cancels, Service B will stop, but **Service C will keep running to completion.** This wastes CPU and database resources on a response that will just be thrown away. |
| **Silent Deadlock/Hangs** | Service C will use its own default timeout (or none at all). If Service C hangs, Service B might hold up its resources indefinitely, risking a cascading failure across your cluster. |

> **Pro Tip for the Interview:** Mention that to prevent developers from forgetting this, production microservices use **gRPC Interceptors** (middleware) to automatically inject, extract, and propagate tracing/deadlines globally without relying on manual code inside every single endpoint handler.