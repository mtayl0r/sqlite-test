package main

import (
    "os"
    "time"
    "flag"
    "fmt"
    "log"
    "strconv"
    "math/rand"
    "sync"
    "net/http"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
    _ "github.com/lib/pq"
)


var insertCount = 0

var insertTxt string
var db *sql.DB
var insertChannel = make(chan string)
var mutex = &sync.Mutex{}

func main() {
    // parse args
    namePtr := flag.String("name", "foo.db", "Sqlite database file name.")
    syncPtr := flag.String("sync", "normal", "PRAGMA synchronous = sync.")
    jmPtr := flag.String("jm", "wal", "PRAGMA journal_mode = jm.")
    flag.Parse()

    fmt.Printf("sqlite: %s %s %s\n", *namePtr, *syncPtr, *jmPtr)

    // setup log file
    f, err := os.OpenFile("channelTest.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()
    log.SetOutput(f)

    // open sqlite connection
    db, _ = sql.Open("sqlite3", *namePtr)
    _, err = db.Exec(fmt.Sprintf("PRAGMA synchronous = %s;", *syncPtr))
    if err != nil {
        panic(err)
    }

    _, err = db.Exec(fmt.Sprintf("PRAGMA journal_mode = %s;", *jmPtr))
    if err != nil {
        panic(err)
    }

    insertTxt = "insert into foo (name) values (?);"
    defer db.Close()

    // set handlers
    http.HandleFunc("/reset", reset)
    http.HandleFunc("/insert", insert)
    http.HandleFunc("/insertMutex", insertMutex)
    http.HandleFunc("/insertGoroutine", insertGoroutine)
    http.ListenAndServe(":8000", nil)
}

func reset(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "reset count: %d", insertCount)
    insertCount = 0
}

func insertMutex(w http.ResponseWriter, r *http.Request) {
    var s = randString(10)
    id := strconv.Itoa(insertCount)
    insertCount++
    s = s + id

    mutex.Lock()
    _, err := db.Exec("insert into foo(name) values('" + s + "');")
    mutex.Unlock()

    if err != nil {
        log.Printf("%s %d\n", err, insertCount)
    }
}

func insertGoroutine(w http.ResponseWriter, r *http.Request) {
    var s = randString(10)
    go insertString(s, insertChannel)
    <-insertChannel
}

func insert(w http.ResponseWriter, r *http.Request) {
    var s = randString(10)
    id := strconv.Itoa(insertCount)
    insertCount++
    s = s + id
    log.Println("start " + id)

    _, err := db.Exec("insert into foo(name) values('" + s + "');")

    if err != nil && err.Error() == "database is locked" {
        log.Printf("%s %d\n", err, insertCount)
        // sleep and retry
        time.Sleep(5 * time.Millisecond)
        _, err = db.Exec("insert into foo(name) values('" + s + "');")
        if err != nil {
            log.Println("--> %s %d\n", err, insertCount)
            // sleep and retry again
            time.Sleep(5 * time.Millisecond)
            _, err = db.Exec("insert into foo(name) values('" + s + "');")
            if err != nil {
                log.Println("----> %s %d\n", err, insertCount)
            }
        }
    }
}

func insertString(s string, c chan string) {
    id := strconv.Itoa(insertCount)
    insertCount++
    s = s + id
    log.Println("start " + id)

    _, err := db.Exec("insert into foo(name) values('" + s + "');")

    if err != nil && err.Error() == "database is locked" {
        log.Printf("%s %d\n", err, insertCount)
    }
    c <- "complete"
}

func randString(len int) string {
    var s = make([]byte, len)
    for i := 0; i < len; i++ {
        s[i] = (byte)(rand.Intn(26)+65)
    }
    return string(s)
}
