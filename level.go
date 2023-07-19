package logger

import "strings"

type Level int

const (
	LevelDebug = Level(iota)
	LevelInfo
	LevelWarn
	LevelError
)

func NewLevel(name string) Level {
	switch strings.ToLower(name) {
	case "dbg", "debug":
		return LevelDebug
	case "inf", "info", "information":
		return LevelInfo
	case "wrn", "warn", "warning":
		return LevelWarn
	case "err", "error":
		return LevelError
	}
	return LevelDebug // default
}

func (l Level) String() string {
	switch l {
	case LevelDebug:
		return "debug"
	case LevelInfo:
		return "info"
	case LevelWarn:
		return "warn"
	case LevelError:
		return "error"
	}
	return ""
}

func (l Level) Prefix() string {
	switch l {
	case LevelDebug:
		return "\x1b[3;30mDBG\x1b[0m"
	case LevelInfo:
		return "\x1b[3;32mINF\x1b[0m"
	case LevelWarn:
		return "\x1b[3;33mWRN\x1b[0m"
	case LevelError:
		return "\x1b[3;34mERR\x1b[0m"
	}
	return ""
}
