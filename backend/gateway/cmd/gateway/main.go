package gateway

import (
	"flag"
	"gateway/internal/app"
	"gateway/internal/config"
	"gateway/internal/infrastructure/logger"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	configFile := flag.String("config", "./configs/config.yaml", "configuration file path")
	flag.Parse()

	cfg, err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	zapLogger := logger.New(&cfg.Logger)

	application := app.New(cfg,
		app.WithLogger(zapLogger),
	)
	go application.Run()

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done

	application.Stop()
}
