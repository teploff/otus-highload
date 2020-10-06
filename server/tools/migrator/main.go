package main

import (
	"database/sql"
	"flag"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

var (
	dir = flag.String("dir", "./migrations", "directory with migration files")
	dsn = flag.String("dsn", "user:password@tcp(localhost:3306)/db?multiStatements=true", "DSN for MySQL")
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
		log.Fatalf("fail to connect with DB, %v\n", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("fail to close connection with DB, %v\n", err)
		}
	}()

	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("fail to create mysql driver, %v\n", err)
	}

	defer func() {
		if err = driver.Close(); err != nil {
			log.Fatalf("fail to close mysql driver, %v\n", err)
		}
	}()

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+*dir,
		"mysql",
		driver,
	)
	if err != nil {
		log.Fatalf("fail to create migrations, %v\n", err)
	}

	switch command {
	case "up":
		if err = m.Up(); err != nil {
			log.Fatalf("fail to up migrations, %v\n", err)
		}
	case "down":
		if err = m.Down(); err != nil {
			log.Fatalf("fail to down migrations, %v\n", err)
		}
	default:
		log.Fatal("unknown operation of migration")
	}
}
