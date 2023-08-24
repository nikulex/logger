package logger

import (
	"fmt"
)

type Config struct {
	LogLevel   string              `json:"logLevel" yaml:"loglevel"`
	ForceDebug bool                `json:"forceDebug" yaml:"forceDebug"`
	Syslog     SyslogOutConfig     `json:"syslog" yaml:"syslog"`
	File       FileOutConfig       `json:"file" yaml:"file"`
	Clickhouse ClickhouseOutConfig `json:"clickhouse" yaml:"clickhouse"`
}

func (cfg *Config) NewLogger() (*Logger, error) {
	outs := make([]LoggerOut, 0)
	outs = append(outs, NewStdOut(&StdOutConfig{
		LogLevel:   NewLevel(cfg.LogLevel),
		ForceDebug: cfg.ForceDebug,
	}))
	if cfg.Syslog.Enabled {
		syslog, err := NewSyslogOut(&cfg.Syslog)
		if err != nil {
			return nil, fmt.Errorf("init syslog error: %w", err)
		}
		outs = append(outs, syslog)
	}
	if cfg.File.Enabled {
		file, err := NewFileOut(&cfg.File)
		if err != nil {
			return nil, fmt.Errorf("init file error: %w", err)
		}
		outs = append(outs, file)
	}
	if cfg.Clickhouse.Enabled {
		clickhouse, err := NewClickhouseOut(&cfg.Clickhouse)
		if err != nil {
			return nil, fmt.Errorf("clickhouse module error: %w", err)
		}
		outs = append(outs, clickhouse)
	}
	return NewLogger(outs...), nil
}

var DefaultConfigMinimal = &Config{
	LogLevel:   "info",
	ForceDebug: false,
}

var DefaultConfig = &Config{
	LogLevel:   "info",
	ForceDebug: false,
	Clickhouse: DefaultClickhouseConfig,
}
