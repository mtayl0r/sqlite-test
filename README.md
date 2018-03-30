# Go Sqlite Insert Test
These are some tests of Sqlite ([using mattn's go-sqlite3 bindings](https://github.com/mattn/go-sqlite3)) insert performance following the findings from [this stackoverflow link](https://stackoverflow.com/questions/1711631/improve-insert-per-second-performance-of-sqlite).

Inserts are tested by inserting random strings against a table:

```sql
CREATE TABLE FOO (
    ID PRIMARY KEY,
    NAME VARCHAR(20)
);
```

## simpleTest.go
Performs serial inserts (ie one after the other), no concurrency:
1. In a transaction.
2. As a prepared statement.
3. Plain insert.

User an specify sqlite or postgres via cmdline args.

## httpTest.go
Simple http server that performs an insert per request. Uses either postgres or sqlite (via cmdline args). This is to test concurrent inserts (via testing with Apache ab for example).
