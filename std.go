package logger

import (
	"log"
	"os"
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
	out *log.Logger
	err *log.Logger
	cfg *StdOutConfig
}

func NewStdOut(cfg *StdOutConfig) *StdOut {
	if cfg == nil {
		cfg = DefaultStdOutConfig
	}
	return &StdOut{
		out: log.New(os.Stdout, "", log.LstdFlags),
		err: log.New(os.Stderr, "", log.LstdFlags),
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
		if level == LevelError {
			l.err.Print(format(level, true, s, i))
		} else {
			l.out.Print(format(level, true, s, i))
		}
	}
}
