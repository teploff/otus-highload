package main

import (
	"context"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	cfg := replication.BinlogSyncerConfig{
		ServerID: 2,
		Flavor:   "mysql",
		Host:     "localhost",
		Port:     3306,
		User:     "replica",
		Password: "oTUSlave#2020",
	}
	syncer := replication.NewBinlogSyncer(cfg)

	binlogFile := "mysql-bin.000001"
	var binlogPos uint32 = 889
	// Start sync with specified binlog file and position
	streamer, err := syncer.StartSync(mysql.Position{Name: binlogFile, Pos: binlogPos})
	if err != nil {
		logger.Fatal("failed start sync", zap.Error(err))
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func(ctx context.Context) {
		for {
			ev, err := streamer.GetEvent(ctx)
			if err != nil {
				logger.Fatal("failed get the binlog event", zap.Error(err))
			}

			select {
			case <-ctx.Done():
				logger.Info("shutting down")

				return
			default:
				parse(ev, logger)
				////logger.Info(fmt.Sprintf("event type = %s", ev.Header.EventType))
				////ev.Header.Dump(os.Stdout)
				////ev.Event.Dump(os.Stdout)
				//ev.Dump(os.Stdout)
			}
		}

	}(ctx)

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done
	cancel()
	<-time.After(time.Second * 5)
}

func parse(event *replication.BinlogEvent, logger *zap.Logger) error {
	switch event.Header.EventType {
	case replication.QUERY_EVENT:
		v, ok := event.Event.(*replication.QueryEvent)
		if !ok {
			return fmt.Errorf("unknow row event")
		}

		logger.Info(fmt.Sprintf("EventType[%s], Schema[%s], Query[%s]", event.Header.EventType, v.Schema, v.Query))
	case replication.TABLE_MAP_EVENT:
		v, ok := event.Event.(*replication.TableMapEvent)
		if !ok {
			return fmt.Errorf("unknow table map event")
		}

		logger.Info(fmt.Sprintf("EventType[%s], Schema[%s], TableID[%d], Table[%s]", event.Header.EventType, v.Schema, v.TableID, v.Table))
	case replication.WRITE_ROWS_EVENTv2:
		v, ok := event.Event.(*replication.RowsEvent)
		if !ok {
			return fmt.Errorf("unknow rows event")
		}

		logger.Info(fmt.Sprintf("EventType[%s], TableID[%d], Rows[%v]", event.Header.EventType, v.TableID, v.Rows))
	}
	return nil
}
