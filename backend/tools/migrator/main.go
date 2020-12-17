package main

import (
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
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

	if err = retryConn(db, 5); err != nil {
		log.Fatalf("fail to connect with DB, %v\n", err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("fail to close connection with DB, %v\n", err)
		}
	}()

	d, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		log.Fatalf("fail to create mysql driver, %v\n", err)
	}

	defer func() {
		if err = d.Close(); err != nil {
			log.Fatalf("fail to close mysql driver, %v\n", err)
		}
	}()

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+*dir,
		"mysql",
		d,
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

func retryConn(db *sql.DB, maxAttempts int) error {
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-ticker.C:
			if i > maxAttempts {
				return fmt.Errorf("attempts to connect with db are over")
			}
			i += 1

			log.Println("try to connect...")

			err := db.Ping()
			switch err {
			case nil:
				log.Println("successful connection")

				return nil
			case driver.ErrBadConn:
				continue
			default:
				return err
			}
		}
	}
}
