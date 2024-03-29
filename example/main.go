package main

import (
	"fmt"
	stdlog "log"
	"time"

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

	// можно добавить параметры рядом с префиксом
	withParams := sublog.Params(logger.Param{"test", 100})
	withParams.Infoln("Hello with params")
	withParams.Params(logger.Param{"some", "hello"}).Infoln("Hello with params 2")

	// можно использовать конкретный модуль
	sublog.Get("clickhouse").Warnln("what happend???")
	subsublog.Std().Debugln("here")

	// чтоб не ждать автоматической выгрузки логов(BatchTime) можно вызвать явно
	log.Flush()
}

func initLogger() *logger.Logger {
	// пример настройки стандартного логгера
	std := logger.NewStdOut(&logger.StdOutConfig{
		Enabled:    true,
		LogLevel:   logger.NewLevel("info"),
		ForceDebug: true, // печатает сообщения с логлевелом debug вне зависимости от настроек логлевела
	})

	// простой файл со всеми левелами
	fileout, err := logger.NewFileOut(&logger.FileOutConfig{
		Enabled:  true,
		FilePath: "myservice.all.log",
		LFlags:   stdlog.LstdFlags,
	})
	if err != nil {
		stdlog.Fatal(err)
	}

	// запись через syslog демон
	syslog, err := logger.NewSyslogOut(&logger.SyslogOutConfig{
		Enabled:  true,
		Facility: "user", // текстовая версия facility в syslog.Priority
		Tag:      "myservice",
	})
	if err != nil {
		stdlog.Fatal(err)
	}

	// в модуле clickhouse используется stdout модуль через уже существующий логгер
	clickhouse, err := logger.NewClickhouseOut(&logger.ClickhouseOutConfig{
		Enabled:        true,
		ClickhouseAddr: "tcp://localhost:9000",
		Database:       "default",
		Username:       "",
		Password:       "",
		Service:        "myservice",
		Timeout:        10 * time.Second,
		BatchTime:      30 * time.Second, // интервал записи в базу
		BatchBuffer:    10000,            // размер буфера логов для записи, при переполнении новые логи игнорируются
	})
	if err != nil {
		stdlog.Fatalf("clickhouse module error: %v\n", err)
	}

	// variadic список модулей в конструкторе
	// если нет stdout модуля, то создатся стандартный
	return logger.NewLogger(std, fileout, syslog, clickhouse)
}
