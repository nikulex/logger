package logger

import (
	"log"
)

type StdOutConfig struct {
	Enabled    bool  `json:"enabled" yaml:"enabled"`
	LogLevel   Level `json:"logLevel" yaml:"logLevel"`
	ForceDebug bool  `json:"forceDebug" yaml:"forceDebug"`
}

var DefaultStdOutConfig *StdOutConfig = &StdOutConfig{
	Enabled:    true,
	LogLevel:   LevelError,
	ForceDebug: false,
}

type StdOut struct {
	l   *log.Logger
	cfg *StdOutConfig
}

func NewStdOut(l *log.Logger, cfg *StdOutConfig) *StdOut {
	if l == nil {
		l = log.Default()
	}
	if cfg == nil {
		cfg = DefaultStdOutConfig
	}
	return &StdOut{
		l:   l,
		cfg: cfg,
	}
}

func (l *StdOut) Close() error {
	return nil
}

func (l *StdOut) init(_ *Logger) {
}

func (l *StdOut) name() string {
	return "std"
}

func (l *StdOut) flush() {
}

func (l *StdOut) log(level Level, s string, i *info) {
	if level >= l.cfg.LogLevel || (level == LevelDebug && l.cfg.ForceDebug) {
		l.l.Print(format(level, s, i))
	}
}
