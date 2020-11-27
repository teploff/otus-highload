package implementation

import (
	"context"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"go.uber.org/zap"
	"replicator/internal/config"
)

type mySQLSyncer struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer

	cancelFunc context.CancelFunc
	doneCh     chan struct{}

	logger *zap.Logger
}

func NewMySQLSyncer(config *config.Config, logger *zap.Logger) (*mySQLSyncer, error) {
	cfg := replication.BinlogSyncerConfig{
		ServerID: config.SlaveID,
		Flavor:   "mysql",
		Host:     config.MySQL.Host,
		Port:     config.MySQL.Port,
		User:     config.MySQL.User,
		Password: config.MySQL.Password,
	}
	syncer := replication.NewBinlogSyncer(cfg)

	// Start sync with specified binlog file and position
	streamer, err := syncer.StartSync(mysql.Position{Name: config.BinlogFile, Pos: config.BinlogPos})
	if err != nil {
		return nil, err
	}

	return &mySQLSyncer{
		syncer:   syncer,
		streamer: streamer,
		doneCh:   make(chan struct{}),
		logger:   logger,
	}, nil
}

func (s *mySQLSyncer) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancelFunc = cancel

	for {
		ev, err := s.streamer.GetEvent(ctx)
		if err != nil {
			s.logger.Error("failed get the binlog event", zap.Error(err))
		}

		select {
		case <-ctx.Done():
			s.logger.Info("mysql syncer binlog shutting down")

			close(s.doneCh)

			return
		default:
			parse(ev, s.logger)
			////logger.Info(fmt.Sprintf("event type = %s", ev.Header.EventType))
			////ev.Header.Dump(os.Stdout)
			////ev.Event.Dump(os.Stdout)
			//ev.Dump(os.Stdout)
		}
	}
}

func (s *mySQLSyncer) Stop() {
	s.cancelFunc()

	<-s.doneCh
	s.syncer.Close()
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
