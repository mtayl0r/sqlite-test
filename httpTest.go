package main

import (
    "fmt"
    "strconv"
    "math/rand"
    "net/http"
    "database/sql"
    //_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

func simple(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hellow world")
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
	_, err = db.Exec("insert into foo(name) values('" + s + "');")
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

//var db, err = sql.Open("sqlite3", "./foo.db")
var db, err = sql.Open("postgres", "user=mark password=00zerozero dbname=test2 sslmode=disable")
var insertCount = 0

func main() {
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
