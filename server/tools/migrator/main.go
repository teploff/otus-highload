package main

import (
	"database/sql"
	"flag"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
)

const mainCountArgs = 3

var (
	dir = flag.String("dir", ".", "directory with migration files")
	dsn = flag.String("dsn", "user:password@tcp(localhost:3306)/db", "DSN for MySQL")
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) < 1 {
		flag.Usage()
		return
	}

	command := args[0]
	db, err := sql.Open("mysql", *dsn)

	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	var arguments []string
	if len(args) > mainCountArgs {
		arguments = append(arguments, args[mainCountArgs:]...)
		log.Println(args)
		log.Println(arguments)
	}

	if err = goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose %v: %v", command, err)
	}
}
