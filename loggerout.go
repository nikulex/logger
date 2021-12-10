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

func format(l Level, s string, i *info) string {
	var prefix, params, sep string
	if len(i.prefix) > 0 {
		prefix = fmt.Sprintf("(%s)", i.prefix)
	}
	if len(i.params) > 0 {
		params = i.params.Json()
	}
	if len(prefix)+len(params) > 0 {
		sep = ":"
	}
	return fmt.Sprintf("[%s]%s%v%s %s", l.Prefix(), prefix, params, sep, s)
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
