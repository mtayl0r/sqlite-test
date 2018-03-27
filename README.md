# Go Sqlite Insert Test
These are some tests of Sqlite [using mattn's go-sqlite3 bindings](https://github.com/mattn/go-sqlite3) insert performance following the findings from [this stackoverflow link](https://stackoverflow.com/questions/1711631/improve-insert-per-second-performance-of-sqlite).

## simpleTest.go
Performs batch inserts (ie one after the other):
1. In a transaction.
2. As a prepared statement.
3. Plain insert.
User an specify sqlite or postgres.

## httpTest.go
Performs insert against either postgres or sqlite (via cmdline args). This is tests concurrent inserts (via testing with Apache ab for example).