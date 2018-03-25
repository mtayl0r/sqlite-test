package main

import (
    "os"
	"flag"
	"fmt"
	"log"
	"time"
    //"math/rand"
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	_ "github.com/lib/pq"
)

func main() {
    sqliteCmd := flag.NewFlagSet("sqlite", flag.ExitOnError)
    namePtr := sqliteCmd.String("name", "foo.db", "Sqlite database file name.")

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

    var insertSql string
    var db *sql.DB

    if postgresCmd.Parsed() {
        connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", *userPtr, *passwordPtr, *dbNamePtr)
        fmt.Printf(connStr + "\n")
	    db, _ = sql.Open("postgres", connStr)
        insertSql = "insert into foo (name) values ($1);"
	    defer db.Close()
    } else if sqliteCmd.Parsed() {
        fmt.Printf("sqlite: %s\n", *namePtr)
	    db, _ = sql.Open("sqlite3", *namePtr)
        insertSql = "insert into foo (name) values (?);"
	    defer db.Close()
    } else {
        fmt.Printf("sqlite:\n")
        sqliteCmd.PrintDefaults()
        fmt.Printf("postgres:\n")
        postgresCmd.PrintDefaults()
        os.Exit(1)
    }

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
        sql := fmt.Sprintf(fmt.Sprintf("insert into foo (name) values ('こんにちわ世界%03d');", i))
		_, err = db.Exec(sql)
		if err != nil {
			log.Fatal(err)
		}
	}

    elapsed = (time.Now().UnixNano() - start) / 1000000
    fmt.Println(elapsed)

    // simple select
    stmt, err = db.Prepare("select count(*) as count from foo;")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var count string
	err = stmt.QueryRow().Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(count)

}

func randString(len int) string {
/*    var s = make([]byte, len)
    for i := 0; i < len; i++ {
        s[i] = (byte)(rand.Intn(26)+65)
    }
    return string(s)
*/
    return "random string"
}

