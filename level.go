package logger

import "strings"

type Level int

const (
	LevelUnknown = Level(iota)
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

func NewLevel(name string) Level {
	switch strings.ToLower(name) {
	case "trace":
		return LevelTrace
	case "dbg", "dbug", "debug":
		return LevelDebug
	case "inf", "info", "information", "informational", "notice":
		return LevelInfo
	case "wrn", "warn", "warning":
		return LevelWarn
	case "err", "eror", "error":
		return LevelError
	case "emerg", "fatal", "alert", "crit", "critical":
		return LevelFatal
	}
	return LevelDebug // default
}

func (l Level) String() string {
	switch l {
	case LevelTrace:
		return "trace"
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	case LevelFatal:
		return "fatal"
	}
	return ""
}

func (l Level) Prefix() string {
	switch l {
	case LevelTrace:
		return "TRACE"
	case LevelDebug:
		return "DEBUG"
	case LevelInfo:
		return "INFO"
	case LevelWarn:
		return "WARN"
	case LevelError:
		return "ERROR"
	case LevelFatal:
		return "FATAL"
	}
	return ""
}

func (l Level) PrefixColor() string {
	switch l {
	case LevelTrace:
		return "\x1b[1;36mTRACE\x1b[0m"
	case LevelDebug:
		return "\x1b[1;34mDEBUG\x1b[0m"
	case LevelInfo:
		return "\x1b[1;32mINFO\x1b[0m"
	case LevelWarn:
		return "\x1b[1;33mWARN\x1b[0m"
	case LevelError:
		return "\x1b[1;31mERROR\x1b[0m"
	case LevelFatal:
		return "\x1b[1;35mFATAL\x1b[0m"
	}
	return ""
}
