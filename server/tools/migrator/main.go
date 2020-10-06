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

const mainCountArgs = 3

var (
	dir = flag.String("dir", "./migrations", "directory with migration files")
	dsn = flag.String("dsn", "user:password@tcp(localhost:3306)/db?multiStatements=true", "DSN for MySQL")
)

func main() {
	flag.Parse()
	//args := flag.Args()

	db, _ := sql.Open("mysql", "user:password@tcp(localhost:3306)/social-network?multiStatements=true")
	driver, _ := mysql.WithInstance(db, &mysql.Config{})
	m, _ := migrate.NewWithDatabaseInstance(
		"file://./migrations",
		"mysql",
		driver,
	)

	if err := m.Down(); err != nil {
		log.Fatal(err)
	}
	//if err:=m.Steps(2); err != nil {
	//	log.Fatal(err)
	//}

	//if len(args) < 1 {
	//	flag.Usage()
	//	return
	//}
	//
	//command := args[0]
	////db, err := sql.Open("mysql", *dsn)
	////if err != nil {
	////	log.Fatalf("goose: failed to open DB: %v\n", err)
	////}
	////defer db.Close()
	//
	//p := &mysql.Mysql{}
	//db, err := p.Open(*dsn)
	//if err != nil {
	//	log.Fatalf("migration: failed to open DB: %v\n", err)
	//}
	//
	//defer func() {
	//	if err = db.Close(); err != nil {
	//		log.Fatalf("migration: failed to close DB: %v\n", err)
	//	}
	//}()
	//
	////
	////
	////
	////if err = db.Ping(); err != nil {
	////	log.Fatalf("mysql ping fail, %v\n", err)
	////}
	//
	//var arguments []string
	//if len(args) > mainCountArgs {
	//	arguments = append(arguments, args[mainCountArgs:]...)
	//}
	//
	//m, err := migrate.NewWithDatabaseInstance("file://"+*dir, "social-network", db)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//switch command {
	//case "up":
	//	m.Up()
	//case "down":
	//	m.Down()
	//}
	//if err = goose.Run(command, db, *dir, arguments...); err != nil {
	//	log.Fatalf("goose %v: %v", command, err)
	//}
}
