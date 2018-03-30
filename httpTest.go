package main

import (
    "os"
    "flag"
    "fmt"
    "strconv"
    "math/rand"
    "net/http"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/lib/pq"
)

func simple(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello world")
    insertCount = 0
}

func query(w http.ResponseWriter, r *http.Request) {
    var s = querySql()
    fmt.Fprintf(w, "Hi there, I love %s!", s)
}

func querySql() string {
    rows, err := db.Query("select * from foo where id = 2;")
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    var name string
    for rows.Next() {
        var id int
        err = rows.Scan(&id, &name)
        if err != nil {
            panic(err)
        }
    }
    return name
}

func insertSql(s string) int {
    insertCount += 1
    s = s + strconv.Itoa(insertCount)
    _, err := db.Exec("insert into foo(name) values('" + s + "');")
    if err != nil {
        return 0
    }
    return 1
}

func mix(w http.ResponseWriter, r *http.Request) {
    var i = rand.Intn(10)
    if i > 6 {
        // insert
        var s = randString(10)
        var i = insertSql(s)
        fmt.Fprintf(w, "insert: %i %s", i, s)
    } else {
        // query
        var s = querySql()
        fmt.Fprintf(w, "query: %s", s)
    }
}

func insert(w http.ResponseWriter, r *http.Request) {
    var s = randString(10)
    var i = insertSql(s)
    fmt.Fprintf(w, "insert: %i %s", i, s)
}

var insertCount = 0

var insertTxt string
var db *sql.DB

func main() {

    sqliteCmd := flag.NewFlagSet("sqlite", flag.ExitOnError)
    namePtr := sqliteCmd.String("name", "foo.db", "Sqlite database file name.")
    syncPtr := sqliteCmd.String("sync", "full", "PRAGMA synchronous = sync.")
    jmPtr := sqliteCmd.String("jm", "delete", "PRAGMA journal_mode = jm.")

    postgresCmd := flag.NewFlagSet("postgres", flag.ExitOnError)
    userPtr := postgresCmd.String("username", "", "Postgres username. (Required)")
    passwordPtr := postgresCmd.String("password", "", "Postgres password. (Required)")
    dbNamePtr := postgresCmd.String("dbName", "", "Postgres database name. (Required)")

    switch os.Args[1] {
        case "sqlite":
            sqliteCmd.Parse(os.Args[2:])
        case "postgres":
            postgresCmd.Parse(os.Args[2:])
    }

    if postgresCmd.Parsed() {
        connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", *userPtr, *passwordPtr, *dbNamePtr)
        fmt.Printf(connStr + "\n")
        db, _ = sql.Open("postgres", connStr)
        insertTxt = "insert into foo (name) values ($1);"
        defer db.Close()
    } else if sqliteCmd.Parsed() {
        fmt.Printf("sqlite: %s %s %s\n", *namePtr, *syncPtr, *jmPtr)
        db, _ = sql.Open("sqlite3", *namePtr)

        _, err := db.Exec(fmt.Sprintf("PRAGMA synchronous = %s;", *syncPtr))
        if err != nil {
            panic(err)
        }

        _, err = db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", *jmPtr))
        if err != nil {
            panic(err)
        }

        insertTxt = "insert into foo (name) values (?);"
        defer db.Close()
    } else {
        fmt.Printf("sqlite:\n")
        sqliteCmd.PrintDefaults()
        fmt.Printf("postgres:\n")
        postgresCmd.PrintDefaults()
        os.Exit(1)
    }


    http.HandleFunc("/simple", simple)
    http.HandleFunc("/query", query)
    http.HandleFunc("/insert", insert)
    http.HandleFunc("/mix", mix)
    http.ListenAndServe(":8000", nil)
}

func randString(len int) string {
    var s = make([]byte, len)
    for i := 0; i < len; i++ {
        s[i] = (byte)(rand.Intn(26)+65)
    }
    return string(s)
}
