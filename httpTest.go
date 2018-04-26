package main

import (
    "os"
    "flag"
    "fmt"
    "sync"
    "math/rand"
    "net/http"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/lib/pq"
)

func insert(w http.ResponseWriter, r *http.Request) {
    var s = randString(10)

    _, err := db.Exec("insert into foo(name) values('" + s + "');")

    if err != nil {
        fmt.Fprintf(w, "insert error")
        return
    }

    fmt.Fprintf(w, "insert: %s", s)
}

var mutex = sync.Mutex{}
func insertMutex(w http.ResponseWriter, r *http.Request) {
     var s = randString(10)

    mutex.Lock()
    _, err := db.Exec("insert into foo(name) values('" + s + "');")
    mutex.Unlock()

    if err != nil {
        fmt.Fprintf(w, "insertMutex error")
        return
    }

    fmt.Fprintf(w, "insertMutex: %s", s)
}

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
        defer db.Close()
    } else if sqliteCmd.Parsed() {
        fmt.Printf("sqlite: %s %s %s\n", *namePtr, *syncPtr, *jmPtr)
        db, _ = sql.Open("sqlite3", *namePtr)
        defer db.Close()

        _, err := db.Exec(fmt.Sprintf("PRAGMA synchronous = %s;", *syncPtr))
        if err != nil {
            panic(err)
        }

        _, err = db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", *jmPtr))
        if err != nil {
            panic(err)
        }

    } else {
        fmt.Printf("sqlite:\n")
        sqliteCmd.PrintDefaults()
        fmt.Printf("postgres:\n")
        postgresCmd.PrintDefaults()
        os.Exit(1)
    }


    http.HandleFunc("/insert", insert)
    http.HandleFunc("/insertMutex", insertMutex)
    http.ListenAndServe(":8000", nil)
}

func randString(len int) string {
    var s = make([]byte, len)
    for i := 0; i < len; i++ {
        s[i] = (byte)(rand.Intn(26)+65)
    }
    return string(s)
}
