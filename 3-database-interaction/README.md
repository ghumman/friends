3. Database interaction with both SQL & NoSQL mental model (medium-high)
They explicitly call out RDBMS + NoSQL. Expect:

*“We have a scan result stored in PostgreSQL (JSONB) and also in Bigtable. Design a Go interface to abstract both. Then show how you’d batch-insert 10k results efficiently without OOM.”*

Or a simpler variant:

“Why does Go’s database/sql require Rows.Close()? What happens if you don’t close it?”