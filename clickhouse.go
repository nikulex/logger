package logger

import (
	"context"
	"fmt"
	"net"
	"time"

	_ "github.com/ClickHouse/clickhouse-go"
	"github.com/jmoiron/sqlx"
)

type ClickhouseOutConfig struct {
	Enabled        bool          `json:"enabled" yaml:"enabled"`
	ClickhouseAddr string        `json:"clickhouseAddr" yaml:"clickhouseAddr"`
	Database       string        `json:"database" yaml:"database"`
	Username       string        `json:"username" yaml:"username"`
	Password       string        `json:"password" yaml:"password"`
	Service        string        `json:"service" yaml:"service"`
	Timeout        time.Duration `json:"timeout" yaml:"timeout"`
	BatchTime      time.Duration `json:"batchTime" yaml:"batchTime"`
	BatchBuffer    int           `json:"batchBuffer" yaml:"batchBuffer"`
}

var DefaultClickhouseConfig = ClickhouseOutConfig{
	Enabled:        true,
	ClickhouseAddr: "tcp://localhost:9000",
	Database:       "default",
	Username:       "",
	Password:       "",
	Service:        "",
	Timeout:        10 * time.Second,
	BatchTime:      30 * time.Second,
	BatchBuffer:    10000,
}

type logData struct {
	Service string
	Server  string
	Level   Level
	Prefix  string
	Params  Params
	Message string
	TM      time.Time
}

type ClickhouseOut struct {
	cfg    *ClickhouseOutConfig
	server string
	conn   *sqlx.DB
	data   chan *logData
	std    *BaseLogger
}

const schema = `
CREATE TABLE IF NOT EXISTS logs(
	service String,
	server  String,
	level   Enum8('debug' = 0, 'info' = 1, 'warn' = 2, 'error' = 3),
	prefix  String,
	params  String,
	message String,
	tm      DateTime
) ENGINE = MergeTree()
ORDER BY tm
PARTITION BY toYYYYMMDD(tm);
`

func getServerName() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", fmt.Errorf("get intergaces error: %v", err)
	}

	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("not found ipv4")
}

func NewClickhouseOut(cfg *ClickhouseOutConfig) (*ClickhouseOut, error) {
	if cfg == nil {
		cfg = &DefaultClickhouseConfig
	}

	connstr := fmt.Sprintf("%s?database=%v&write_timeout=%v",
		cfg.ClickhouseAddr, cfg.Database, cfg.Timeout.Seconds())
	conn, err := sqlx.Open("clickhouse", connstr)
	if err != nil {
		return nil, fmt.Errorf("clickhouse connect error: %v", err)
	}

	if _, err = conn.Exec(schema); err != nil {
		return nil, fmt.Errorf("init schema error: %v", err)
	}

	server, err := getServerName()
	if err != nil {
		return nil, fmt.Errorf("get server ip error: %v", err)
	}

	log := &ClickhouseOut{
		cfg:    cfg,
		server: server,
		conn:   conn,
		data:   make(chan *logData, cfg.BatchBuffer),
	}

	go func() {
		for {
			time.Sleep(cfg.BatchTime)
			log.flush()
		}
	}()

	return log, nil
}

func (l *ClickhouseOut) Close() error {
	return l.conn.Close()
}

func (l *ClickhouseOut) init(log *Logger) {
	l.std = log.New("clickhouse logger").Std()
}

func (l *ClickhouseOut) name() string {
	return "clickhouse"
}

func (l *ClickhouseOut) flush() {
	if len(l.data) == 0 {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), l.cfg.Timeout)
	defer cancel()

	tx, err := l.conn.Begin()
	if err != nil {
		l.std.Errorf("start transaction error: %v", err)
		return
	}
	stmt, err := tx.Prepare("INSERT INTO logs (service, server, level, prefix, params, message, tm) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		l.std.Errorf("prepare stmt error: %v", err)
		return
	}
	defer stmt.Close()

	count := 0
	for len(l.data) > 0 {
		data := <-l.data
		if _, err = stmt.ExecContext(ctx, data.Service, data.Server, data.Level.String(), data.Prefix, data.Params.Json(), data.Message, data.TM); err != nil {
			l.std.Errorf("insert error: %v", err)
			// skip error
		} else {
			count++
		}
	}
	if err = tx.Commit(); err != nil {
		l.std.Errorf("commit error: %v", err)
	} else {
		l.std.Debugf("inserted %v logs", count)
	}
}

func (l *ClickhouseOut) log(level Level, s string, i *info) {
	if len(l.data) < l.cfg.BatchBuffer {
		l.data <- &logData{
			Service: l.cfg.Service,
			Server:  l.server,
			Prefix:  i.prefix,
			Params:  i.params,
			Level:   level,
			Message: s,
			TM:      time.Now(),
		}
	} // else skip
}
