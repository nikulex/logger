package logger

import (
	"log"
	"os"
	"path/filepath"
)

type FileOutConfig struct {
	Enabled  bool   `json:"enabled" yaml:"enabled"`
	FilePath string `json:"filePath" yaml:"filePath"`
	LFlags   int    `json:"lflags" yaml:"lflags"`
}

var DefaultFileOutConfig = &FileOutConfig{
	Enabled:  true,
	FilePath: filepath.Base(os.Args[0]) + ".log",
	LFlags:   log.LstdFlags,
}

type FileOut struct {
	file *os.File
	l    *log.Logger
}

func NewFileOut(cfg *FileOutConfig) (*FileOut, error) {
	file, err := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FileOut{
		file: file,
		l:    log.New(file, "", cfg.LFlags),
	}, nil
}

func (l *FileOut) Close() error {
	return l.file.Close()
}

func (l *FileOut) log(level Level, s string, i *info) {
	l.l.Print(format(level, s, i))
}

func (l *FileOut) init(main *Logger) {
}

func (l *FileOut) name() string {
	return "file"
}

func (l *FileOut) flush() {
}
