package clickhouse

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/rtsoftSG/plugin/toolbox/database"
	"go.uber.org/zap"
	"golang.org/x/sync/semaphore"
)

const timeOut = time.Minute * 30

type tableBuffer struct {
	commands []*database.InsertCommand
	sync.Mutex
}

//MaxParallelQuery may be equal max_concurrent_queries in clickhouse configuration.
const MaxParallelQuery = 100

//Storage clickhouse storage.
type Storage struct {
	db *sql.DB

	buffers  map[string]*tableBuffer
	bufferMu sync.RWMutex

	pushTimeout time.Duration
	doneCh      chan struct{}

	logger *zap.Logger

	sem *semaphore.Weighted
}

//NewStorage create new clickhouse storage.
func NewStorage(db *sql.DB, pushTimeout time.Duration, logger *zap.Logger) *Storage {
	return &Storage{
		db:          db,
		logger:      logger,
		buffers:     make(map[string]*tableBuffer),
		sem:         semaphore.NewWeighted(MaxParallelQuery),
		pushTimeout: pushTimeout,
		doneCh:      make(chan struct{}),
	}
}

func (s *Storage) DB() *sql.DB {
	return s.db
}

func (s *Storage) StartBatching() {
	s.logger.Info("clickhouse batching is starting")

	ticker := time.NewTicker(s.pushTimeout)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.pushBatches()
		case <-s.doneCh:
			close(s.doneCh)
			s.logger.Info("clickhouse batching is over")

			return
		}
	}
}

func (s *Storage) Stop() {
	s.doneCh <- struct{}{}
	<-s.doneCh
}

func (s *Storage) pushBatches() {
	s.bufferMu.RLock()
	buffersToPush := make(map[string]*tableBuffer, len(s.buffers))

	for i := range s.buffers {
		buffersToPush[i] = s.buffers[i]
	}
	s.bufferMu.RUnlock()

	for _, tb := range buffersToPush {
		tb.Lock()

		tableBufferLen := len(tb.commands)
		batch := make([]*database.InsertCommand, tableBufferLen)
		copy(batch, tb.commands)
		// nolint:godox
		// todo perform if needed
		tb.commands = make([]*database.InsertCommand, 0)

		tb.Unlock()

		go func(b []*database.InsertCommand) {
			if err := s.executeCommands(b); err != nil {
				s.logger.Error("clickhouse insert error", zap.Error(err))
			}
		}(batch)
	}
}

//Insert insert data into database.
//note: data will be buffered while batchSize limit not be excited.
func (s *Storage) Insert(cmd *database.InsertCommand) {
	var (
		tableBuff *tableBuffer
		tableName = cmd.TableName()
	)

	s.bufferMu.Lock()
	if _, exists := s.buffers[tableName]; !exists {
		s.buffers[tableName] = &tableBuffer{}
	}

	tableBuff = s.buffers[tableName]
	s.bufferMu.Unlock()

	tableBuff.Lock()
	defer tableBuff.Unlock()

	tableBuff.commands = append(tableBuff.commands, cmd)
}

func (s *Storage) executeCommands(commands []*database.InsertCommand) error {
	if len(commands) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	if err := s.sem.Acquire(ctx, 1); err != nil {
		return fmt.Errorf("insert not executed, cannot acquire ch query semaphore: %w", err)
	}
	defer s.sem.Release(1)

	var (
		fieldCount = len(commands[0].Fields())
		tableName  = commands[0].TableName()
	)

	cols := make([]string, fieldCount)
	for i, field := range commands[0].Fields() {
		cols[i] = field.Name()
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// nolint:gosec
	if err = executeTx(
		tx,
		fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", tableName, strings.Join(cols, ","), strings.Repeat("?,", fieldCount)),
		commands,
	); err != nil {
		_ = tx.Rollback()
		return err
	}

	return nil
}

func executeTx(tx *sql.Tx, sql string, commands []*database.InsertCommand) error {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()

	args := make([]interface{}, len(commands[0].Fields()))

	for _, cmd := range commands {
		for i, field := range cmd.Fields() {
			args[i] = field.Value()
		}

		if _, err := stmt.Exec(args...); err != nil {
			return err
		}
	}

	return tx.Commit()
}
