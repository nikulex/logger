package main

import (
	"fmt"
	"runtime/debug"

	"github.com/nikulex/logger"
)

func main() {
	log := initLogger()
	defer log.Close()

	log.Infof("Hello %v!\n", "World")
	log.Error(fmt.Errorf("some error"))

	// иерархия логгеров, проставляет префикс
	sublog := log.New("mymodule")
	sublog.Infoln("hello sublog")
	subsublog := sublog.New("submodule")
	subsublog.Infoln("hello subsublog")

	withParams := sublog.Params(logger.Param{"test", 100})
	withParams.Infoln("Hello with params")
	withParams.Params(logger.Param{"some", "hello"}).Infoln("Hello with params 2")

	subsublog.Std().Debugln("here")
	subsublog.Params(logger.Param{"stack", string(debug.Stack())}).Tranceln("deep debug")

	subsublog.Fatalln("AAAAAAA")
	sublog.Warnln("OH... OK")
}

func initLogger() *logger.Logger {
	std := logger.NewStdOut(&logger.StdOutConfig{
		Enabled:  true,
		LogLevel: logger.NewLevel("trace"),
	})
	return logger.NewLogger(std)
}
