package implementation

import (
	"context"
	"fmt"
	"github.com/siddontang/go-mysql/mysql"
	"github.com/siddontang/go-mysql/replication"
	"go.uber.org/zap"
	"replicator/internal/config"
	"replicator/internal/infrastructure/tarantool"
)

type mySQLSyncer struct {
	syncer   *replication.BinlogSyncer
	streamer *replication.BinlogStreamer

	tableID uint64

	cancelFunc context.CancelFunc
	doneCh     chan struct{}

	tarantoolConn *tarantool.Conn
	logger        *zap.Logger
}

func NewMySQLSyncer(config *config.Config, tarantoolConn *tarantool.Conn, logger *zap.Logger) (*mySQLSyncer, error) {
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
		syncer:        syncer,
		streamer:      streamer,
		doneCh:        make(chan struct{}),
		tarantoolConn: tarantoolConn,
		logger:        logger,
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
			if err := s.parse(ev); err != nil {
				s.logger.Error("failed parsing event", zap.Error(err))
			}
		}
	}
}

func (s *mySQLSyncer) Stop() {
	s.cancelFunc()

	<-s.doneCh
	s.syncer.Close()
}

func (s *mySQLSyncer) parse(event *replication.BinlogEvent) error {
	switch event.Header.EventType {
	case replication.QUERY_EVENT:
		_, ok := event.Event.(*replication.QueryEvent)
		if !ok {
			return fmt.Errorf("unknow row event")
		}

		//s.logger.Info(fmt.Sprintf("EventType[%s], Schema[%s], Query[%s]", event.Header.EventType, v.Schema,
		//	v.Query))
	case replication.TABLE_MAP_EVENT:
		v, ok := event.Event.(*replication.TableMapEvent)
		if !ok {
			return fmt.Errorf("unknow table map event")
		}

		//s.logger.Info(fmt.Sprintf("EventType[%s], Schema[%s], TableID[%d], Table[%s]", event.Header.EventType,
		//	v.Schema, v.TableID, v.Table))
		if string(v.Table) == "user" {
			s.tableID = v.TableID
		}
	case replication.WRITE_ROWS_EVENTv2:
		v, ok := event.Event.(*replication.RowsEvent)
		if !ok {
			return fmt.Errorf("unknow rows event")
		}
		//s.logger.Info(fmt.Sprintf("EventType[%s], TableID[%d], Rows[%v]", event.Header.EventType, v.TableID,
		//	v.Rows))

		if v.TableID == s.tableID {
			v1, _ := v.Rows[0][2].([]byte)
			v2, _ := v.Rows[0][8].([]byte)

			err := s.tarantoolConn.Insert(v.Rows[0][0], v.Rows[0][1], string(v1), v.Rows[0][3],
				v.Rows[0][4], v.Rows[0][5], v.Rows[0][6], v.Rows[0][7], string(v2), v.Rows[0][11])
			if err != nil {
				return err
			}
		}
	}

	return nil
}
