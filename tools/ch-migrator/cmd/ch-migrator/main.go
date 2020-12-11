package main

import (
	"ch-migrator/internal/config"
	"flag"
	"fmt"
	"log"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
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

	m, err := migrate.New(fmt.Sprintf("file://%s", cfg.MigrationsPath), cfg.DSN)
	if err != nil {
		log.Fatal(err)
	}
	defer m.Close()

	switch cfg.Operation {
	case "up":
		if err = m.Up(); err != nil {
			logger.Fatal("fail up operation", zap.Error(err))
		}

		logger.Info("Successfully up operation!")
	case "down":
		if err = m.Down(); err != nil {
			logger.Fatal("fail up operation", zap.Error(err))
		}

		log.Println("Successfully down operation!")
	default:
		log.Fatal(fmt.Sprintf("fail to recognize operation %s", cfg.Operation))
	}
}
