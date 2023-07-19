package logger

import (
	"fmt"

	syslog "github.com/RackSec/srslog"
)

type SyslogOutConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	Facility string `json:"facility" yaml:"facility"`
	Tag      string `json:"tag" yaml:"tag"`
}

var DefaultSyslogOutConfig = &SyslogOutConfig{
	Enabled:  true,
	Facility: "daemon",
	Tag:      "",
}

type SyslogOut struct {
	w   *syslog.Writer
	std *BaseLogger
}

func NewSyslogOut(cfg *SyslogOutConfig) (*SyslogOut, error) {
	w, err := syslog.New(facility(cfg.Facility), cfg.Tag)
	if err != nil {
		return nil, fmt.Errorf("syslog out init error: %w", err)
	}
	return &SyslogOut{
		w: w,
	}, nil
}

func (l *SyslogOut) Close() error {
	return l.w.Close()
}

func (l *SyslogOut) levelFunc(level Level) func(string) error {
	switch level {
	case LevelInfo:
		return l.w.Info
	case LevelWarn:
		return l.w.Warning
	case LevelError:
		return l.w.Err
	case LevelDebug:
		return l.w.Debug
	default:
		l.std.Warnf("unexpected log level: %v", level)
		return l.w.Notice
	}
}

func (l *SyslogOut) log(level Level, s string, i *info) {
	msg := format(level, false, s, i)
	if err := l.levelFunc(level)(msg); err != nil {
		l.std.Errorf("syslog write error: %v", err)
	}
}

func (l *SyslogOut) init(main *Logger) {
	l.std = main.New("syslog").Std()
}

func (l *SyslogOut) name() string {
	return "syslog"
}

func (l *SyslogOut) flush() {
}

func facility(s string) syslog.Priority {
	switch s {
	case "kern", "kernel":
		return syslog.LOG_KERN
	case "user":
		return syslog.LOG_USER
	case "mail":
		return syslog.LOG_MAIL
	case "daemon":
		return syslog.LOG_DAEMON
	case "auth":
		return syslog.LOG_AUTH
	case "syslog":
		return syslog.LOG_SYSLOG
	case "lpr":
		return syslog.LOG_LPR
	case "news":
		return syslog.LOG_NEWS
	case "uucp":
		return syslog.LOG_UUCP
	case "cron":
		return syslog.LOG_CRON
	case "authpriv":
		return syslog.LOG_AUTHPRIV
	case "ftp":
		return syslog.LOG_FTP
	case "local0":
		return syslog.LOG_LOCAL0
	case "local1":
		return syslog.LOG_LOCAL1
	case "local2":
		return syslog.LOG_LOCAL2
	case "local3":
		return syslog.LOG_LOCAL3
	case "local4":
		return syslog.LOG_LOCAL4
	case "local5":
		return syslog.LOG_LOCAL5
	case "local6":
		return syslog.LOG_LOCAL6
	case "local7":
		return syslog.LOG_LOCAL7
	}
	return syslog.Priority(0)
}
