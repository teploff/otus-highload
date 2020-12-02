package main

import (
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"migrator/internal/config"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate"
	"github.com/golang-migrate/migrate/database/mysql"
	_ "github.com/golang-migrate/migrate/source/file"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		logger.Fatal("error reading config file", zap.Error(err))
	}

	db, err := sql.Open("mysql", cfg.DSN)
	if err != nil {
		logger.Fatal("fail to connect with DB", zap.Error(err))
	}

	if err = retryConn(db, 5, logger); err != nil {
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
		"file://"+cfg.MigrationsPath,
		"mysql",
		d,
	)
	if err != nil {
		log.Fatalf("fail to create migrations, %v\n", err)
	}

	switch cfg.Operation {
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

func retryConn(db *sql.DB, maxAttempts int, logger *zap.Logger) error {
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

			logger.Info("try to connect...")

			err := db.Ping()
			switch err {
			case nil:
				logger.Info("successful connection")

				return nil
			case driver.ErrBadConn:
				continue
			default:
				return err
			}
		}
	}
}
