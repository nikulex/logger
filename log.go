package logger

import (
	"fmt"
)

type BaseLogger struct {
	internal
	*info
}

func (l *BaseLogger) Print(a ...any) {
	l.log(LevelUnknown, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Trace(a ...any) {
	l.log(LevelTrace, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Debug(a ...any) {
	l.log(LevelDebug, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Info(a ...any) {
	l.log(LevelInfo, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Warn(a ...any) {
	l.log(LevelWarn, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Error(a ...any) {
	l.log(LevelError, fmt.Sprint(a...), l.info)
}
func (l *BaseLogger) Fatal(a ...any) {
	l.log(LevelFatal, fmt.Sprint(a...), l.info)
}

func (l *BaseLogger) Printf(format string, a ...any) {
	l.log(LevelUnknown, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Tracef(format string, a ...any) {
	l.log(LevelTrace, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Debugf(format string, a ...any) {
	l.log(LevelDebug, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Infof(format string, a ...any) {
	l.log(LevelInfo, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Warnf(format string, a ...any) {
	l.log(LevelWarn, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Errorf(format string, a ...any) {
	l.log(LevelError, fmt.Sprintf(format, a...), l.info)
}
func (l *BaseLogger) Fatalf(format string, a ...any) {
	l.log(LevelFatal, fmt.Sprintf(format, a...), l.info)
}

func (l *BaseLogger) Println(a ...any) {
	l.log(LevelUnknown, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Tranceln(a ...any) {
	l.log(LevelTrace, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Debugln(a ...any) {
	l.log(LevelDebug, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Infoln(a ...any) {
	l.log(LevelInfo, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Warnln(a ...any) {
	l.log(LevelWarn, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Errorln(a ...any) {
	l.log(LevelError, fmt.Sprintln(a...), l.info)
}
func (l *BaseLogger) Fatalln(a ...any) {
	l.log(LevelFatal, fmt.Sprintln(a...), l.info)
}

type loggerOuts map[string]LoggerOut

func (outs loggerOuts) log(l Level, s string, i *info) {
	for _, out := range outs {
		out.log(l, s, i)
	}
}

type Logger struct {
	BaseLogger
	outs loggerOuts
}

func NewLogger(outs ...LoggerOut) *Logger {
	louts := make(loggerOuts)
	main := &Logger{
		BaseLogger: BaseLogger{louts, &info{}},
		outs:       louts,
	}
	for _, out := range outs {
		out.init(main)
		main.outs[out.name()] = out
	}
	// required stdout logger
	if _, ok := main.outs["std"]; !ok {
		out := NewStdOut(nil)
		out.init(main)
		main.outs[out.name()] = out
	}
	return main
}

func (l *Logger) Close() error {
	for _, out := range l.outs {
		if err := out.Close(); err != nil {
			return fmt.Errorf("close main logger error: %w", err)
		}
	}
	return nil
}

// sublogger with prefix
func (l *Logger) New(name string) *Logger {
	if l.prefix != "" {
		name = l.prefix + "/" + name // submodules path
	}
	return &Logger{
		BaseLogger: BaseLogger{
			internal: l.internal,
			info: &info{
				params: l.params,
				prefix: name,
			},
		},
		outs: l.outs,
	}
}

type Param struct {
	Name  string
	Value interface{}
}

func (l *Logger) Params(params ...Param) *Logger {
	paramsMap := make(Params)
	for _, p := range params {
		paramsMap[p.Name] = p.Value
	}
	return l.ParamsMap(paramsMap)
}

// sublogger with params
func (l *Logger) ParamsMap(params map[string]interface{}) *Logger {
	child := &Logger{
		BaseLogger: BaseLogger{
			internal: l.internal,
			info: &info{
				params: make(Params, len(l.params)+len(params)),
				prefix: l.prefix,
			},
		},
		outs: l.outs,
	}
	for k, v := range l.params {
		child.params[k] = v
	}
	for k, v := range params {
		child.params[k] = v
	}
	return child
}

func (l *Logger) Flush() {
	for _, out := range l.outs {
		out.flush()
	}
}

// access to specific module with same interface
func (l *Logger) Get(out string) *BaseLogger {
	if o, ok := l.outs[out]; ok {
		return &BaseLogger{o, l.info}
	}
	return nil
}

// get stdout logger (ensure not nil)
func (l *Logger) Std() *BaseLogger {
	return l.Get("std")
}
