package logger

import (
	"encoding/json"
	"fmt"
	"io"
)

type Params map[string]interface{}

func (params Params) Json() string {
	if len(params) == 0 {
		return ""
	}
	res, _ := json.Marshal(params)
	return string(res)
}

type info struct {
	params Params
	prefix string
}

func format(l Level, colored bool, s string, i *info) string {
	var prefix, params string
	if len(i.prefix) > 0 {
		if colored {
			prefix = fmt.Sprintf("\x1b[0;30m (%s)\x1b[0m", i.prefix)
		} else {
			prefix = fmt.Sprintf(" (%s)", i.prefix)
		}
	}
	if len(i.params) > 0 {
		if colored {
			params = "\x1b[0;36m" + i.params.Json() + "\x1b[0m"
		} else {
			params = i.params.Json()
		}
	}
	var loglevel string
	if colored {
		loglevel = l.PrefixColor()
	} else {
		loglevel = l.Prefix()
	}
	if l != LevelUnknown {
		loglevel = "[" + loglevel + "]"
	}
	if l == LevelInfo || l == LevelWarn {
		loglevel += " "
	}
	return fmt.Sprintf("%s%s%v %s", loglevel, prefix, params, s)
}

type internal interface {
	log(l Level, s string, i *info)
}

// logger module interface
type LoggerOut interface {
	internal
	io.Closer
	init(l *Logger)
	name() string
	flush()
}
