package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	stdlog "log"
	"os"

	"github.com/nikulex/logger"
)

type Config struct {
	Test string `json:"test"`
	// app configs ...
	Logger *logger.Config `json:"logger"`
}

var DefaultConfig = &Config{
	Test:   "hello",
	Logger: logger.DefaultConfigMinimal,
}

func LoadConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Config) Save(path string) error {
	data, err := json.MarshalIndent(cfg, " ", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, data, 0777)
}

const configPath = "example_config.json"

func main() {
	var config *Config
	if _, err := os.Stat(configPath); err != nil {
		config = DefaultConfig
		config.Save(configPath)
		stdlog.Println("saved example config:", configPath)
	} else {
		config, err = LoadConfig(configPath)
		if err != nil {
			stdlog.Fatalf("failed to load config: %v", err)
		}
		stdlog.Println("loaded config:", configPath)
	}

	forceDebug := flag.Bool("debug", false, "force debug mode")
	flag.Parse()

	if *forceDebug { // override config value
		config.Logger.ForceDebug = *forceDebug
	}

	log, err := config.Logger.NewLogger()
	if err != nil {
		stdlog.Fatalf("failed to init logger: %v", err)
	}
	defer log.Close()

	log.Info(config.Test)
	log.Debug("bug here")
}
