package main

import (
	"fmt"
	"log"
	"time"
    "math/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/lib/pq"
)

func main() {

    insertSql := "insert into foo (name) values (?);"
    //insertSql := "insert into foo (name) values ($1);"

	db, err := sql.Open("sqlite3", "./foo.db")
	//db, err := sql.Open("postgres", "user=mark password=00zerozero dbname=test2 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	/* transaction test */

    tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare(insertSql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

    start := time.Now().UnixNano()

	for i := 0; i < 10000; i++ {
		_, err = stmt.Exec(randString(15))
		if err != nil {
			log.Fatal(err)
		}
	}

    elapsed := (time.Now().UnixNano() - start) / 1000000
    fmt.Println(elapsed)

	tx.Commit()
    
	/* prepared stmt test */

    stmt, err = db.Prepare(insertSql)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

    start = time.Now().UnixNano()

	for i := 0; i < 10000; i++ {
		_, err = stmt.Exec(randString(15))
		if err != nil {
			log.Fatal(err)
		}
	}

    elapsed = (time.Now().UnixNano() - start) / 1000000
    fmt.Println(elapsed)

	/* direct insert stmt test */

    start = time.Now().UnixNano()

	for i := 0; i < 10000; i++ {
        sql := fmt.Sprintf("insert into foo (name) values ('こんにちわ世界%03d');", i)
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}

    elapsed = (time.Now().UnixNano() - start) / 1000000
    fmt.Println(elapsed)

    // simple select
    stmt, err = db.Prepare("select name from foo where id = ?")
    //stmt, err = db.Prepare("select name from foo where id = $1")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var name string
	err = stmt.QueryRow("1").Scan(&name)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(name)

}

func randString(len int) string {
    var s = make([]byte, len)
    for i := 0; i < len; i++ {
        s[i] = (byte)(rand.Intn(26)+65)
    }
    return string(s)
}

