package internal

import (
	"context"
	"database/sql"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

type connection struct {
	log     *zap.Logger
	ctx     context.Context
	lock    chan struct{}
	unlock  chan struct{}
	mux     sync.Mutex
	db      *sql.DB
	Name    string
	timeout int
}

type ConnManager interface {
	UnlockAll() error
	Connect(string, []string, int) error
	LockAll() error
}

type connManager struct {
	ctx         context.Context
	log         *zap.Logger
	Connections map[string]*connection
}

type ConnManagerConfig *ConnManagerCfg

type ConnManagerCfg struct {
	LockTimeout int
	Path        string
	DBs         []string
}

func NewConnManagerConfig() ConnManagerConfig {
	return &ConnManagerCfg{}
}

func ProvideConnManager(ctx context.Context, cfg ConnManagerConfig, logger *zap.Logger) (ConnManager, error) {
	connections := make(map[string]*connection)

	connman := &connManager{
		ctx:         ctx,
		log:         logger,
		Connections: connections,
	}

	if err := connman.Connect(cfg.Path, cfg.DBs, cfg.LockTimeout); err != nil {
		return nil, err
	}
	return connman, nil
}

func (c *connManager) Connect(path string, databases []string, lockTimeout int) error {
	for _, database := range databases {
		fullPath := filepath.Join(path, database)
		c.log.Info("Connecting to database", zap.String("database", fullPath))
		conn, _ := sql.Open("sqlite3", fullPath)
		if err := conn.Ping(); err != nil {
			c.log.Error("Error connecting to database", zap.String("database", fullPath))
			return err
		}
		c.Connections[database] = &connection{
			lock:    make(chan struct{}),
			unlock:  make(chan struct{}),
			timeout: lockTimeout,
			Name:    database,
			log:     c.log,
			ctx:     c.ctx,
			db:      conn,
		}
	}

	return nil
}

func (c *connection) Lock(wg *sync.WaitGroup, conn ConnManager) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case <-c.unlock:
			if c.mux.TryLock() {
				c.log.Info("Database not locked", zap.String("database", c.Name))
				c.mux.Unlock()
				continue
			}
			c.log.Info("Unlocking by user request", zap.String("database", c.Name))
			if _, err := c.db.ExecContext(c.ctx, "COMMIT;"); err != nil {
				c.log.Error("Error unlocking database", zap.String("database", c.Name), zap.Error(err))
			}
			c.mux.Unlock()
			return
		case <-time.After(time.Duration(c.timeout) * time.Second):
			if c.mux.TryLock() {
				c.log.Info("Database not locked", zap.String("database", c.Name))
				c.mux.Unlock()
				continue
			}
			c.log.Info("Unlocking by timeout", zap.String("database", c.Name))
			if _, err := c.db.ExecContext(c.ctx, "COMMIT;"); err != nil {
				c.log.Error("Error unlocking database", zap.String("database", c.Name), zap.Error(err))
			}
			c.mux.Unlock()
			return
		case <-c.lock:
			var i int
			for i = 0; i < 3; i++ {
				if _, err := c.db.ExecContext(c.ctx, "BEGIN EXCLUSIVE TRANSACTION;"); err != nil {
					c.log.Error("Error locking database", zap.String("database", c.Name), zap.Error(err))
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}
			if i == 3 {
				c.log.Error("Error; Unlocking other databases", zap.String("database", c.Name))
				c.mux.Unlock()
				conn.UnlockAll()

			}
			c.log.Info("Database locked", zap.String("database", c.Name))
		}
	}
}

func (c *connManager) LockAll() error {
	go func() {
		var wg sync.WaitGroup
		for _, conn := range c.Connections {
			if conn.mux.TryLock() {
				go conn.Lock(&wg, c)
				conn.lock <- struct{}{}
			}
		}

		wg.Wait()
	}()

	return nil
}

func (c *connManager) UnlockAll() error {
	c.log.Info("Unlocking all databases")
	for _, conn := range c.Connections {
		if conn.mux.TryLock() {
			c.log.Info("Database not locked", zap.String("database", conn.Name))
			conn.mux.Unlock()
			continue
		}
		conn.unlock <- struct{}{}
	}

	return nil
}
