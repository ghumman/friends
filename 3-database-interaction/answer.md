Here is your preparation guide for the database section. Given Synack’s focus on long-running security scans, efficiency and resource management are paramount.

---

## 1. Abstracting PostgreSQL & Bigtable in Go

To abstract both databases, you need an interface that handles data generically. Since PostgreSQL uses `JSONB` and Bigtable works with unstructured column-family bytes, passing standard Go structs or marshaled slices is the cleanest approach.

### The Interface Design

```go
package storage

import "context"

// ScanResult represents the domain model for a security scan
type ScanResult struct {
    ID        string
    Target    string
    Findings  []byte // Raw JSON or serialized payload
}

// ScanRepository abstracts the underlying database implementation
type ScanRepository interface {
    // BatchInsert stores a chunk of results efficiently
    BatchInsert(ctx context.Context, results []ScanResult) error
}

```

### Implementing 10k Batch Inserts Without Out-of-Memory (OOM)

Loading 10,000 heavy scan results into memory all at once can spike your heap and trigger the Linux OOM killer. The senior approach is to stream the data or process it in **micro-batches** (e.g., chunks of 500 or 1,000) using a worker pool or sequential pipeline.

Here is how you implement the batching logic for both mental models:

#### Mental Model A: PostgreSQL (Relational/JSONB)

Instead of executing 10,000 individual `INSERT` statements (which destroys performance due to network round-trips), you combine multiple rows into a single query or use the `COPY` protocol.

```go
func (r *PostgresRepo) BatchInsert(ctx context.Context, results []ScanResult) error {
    // Chunking 10k results into sizes of 1000 to keep memory footprint flat
    const chunkSize = 1000
    
    for i := 0; i < len(results); i += chunkSize {
        end := i + chunkSize
        if end > len(results) {
            end = len(results)
        }
        chunk := results[i:end]

        // Build a single bulk INSERT statement: 
        // INSERT INTO scans (id, target, data) VALUES ($1, $2, $3), ($4, $5, $6)...
        // Or use pgx.Conn.CopyFrom for maximum performance.
        if err := r.executeBulkInsert(ctx, chunk); err != nil {
            return err
        }
    }
    return nil
}

```

#### Mental Model B: Bigtable (NoSQL/Key-Value)

Bigtable excels at massive scale, but individual mutations are limited. You use the cloud SDK's built-in **Bulk Mutator** API, which automatically batches rows and flushes them based on size or time boundaries.

```go
func (r *BigtableRepo) BatchInsert(ctx context.Context, results []ScanResult) error {
    tbl := r.client.Open(r.tableName)
    
    // ApplyBulk accepts a slice of keys and a slice of mutations
    muts := make([]*bigtable.Mutation, len(results))
    rowKeys := make([]string, len(results))
    
    for i, res := range results {
        rowKeys[i] = res.ID
        mut := bigtable.NewMutation()
        mut.Set("info", "data", bigtable.Now(), res.Findings)
        muts[i] = mut
    }
    
    // ApplyBulk handles the parallel network distribution natively
    errs, err := tbl.ApplyBulk(ctx, rowKeys, muts)
    if err != nil {
        return err
    }
    // Check individual row errors inside 'errs' slice...
    return nil
}

```

---

## 2. Why does Go's `database/sql` require `Rows.Close()`?

This is a classic "gotcha" question that tests if you understand what Go is doing under the hood with connection pooling.

### What happens conceptually?

When you run a query using `db.Query()`, Go fetches a database connection out of its active connection pool (`sql.DB`) and locks it to your current operation.

As you call `Rows.Next()`, Go streams the data from the database server over that specific network connection.

### What happens if you *don't* close it?

If you forget to call `Rows.Close()`, or if your function exits early due to an error before reaching the end of the data stream, **the database connection is never released back to the pool.** It remains permanently locked and "in-use" by that abandoned `Rows` object.

### The Consequences:

1. **Connection Pool Exhaustion:** Over time, as your application handles more traffic, more connections get locked open. Eventually, the pool hits `SetMaxOpenConns`. New incoming requests will block forever waiting for a free connection, causing the service to freeze or timeout.
2. **Database Resource Leaks:** On the database server itself (like PostgreSQL), those orphaned connections remain active backends. They consume server memory, file descriptors, and locks, eventually degrading database performance or hitting the database’s `max_connections` ceiling.

### The Idiomatic Go Fix

Always use `defer rows.Close()` immediately after checking the query error:

```go
rows, err := db.QueryContext(ctx, "SELECT id, data FROM scans")
if err != nil {
    return err
}
defer rows.Close() // Guarantees connection release even if code panics or exits early

for rows.Next() {
    // scan data...
}

// Always check rows.Err() after the loop finishes to catch streaming errors
if err = rows.Err(); err != nil {
    return err
}

```