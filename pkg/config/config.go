package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"timeMachine/delivery/grpcserver"
	"timeMachine/delivery/httpserver"
	"timeMachine/repository/postgres"
	"timeMachine/scheduler"
)

type Config struct {
	HTTPServer httpserver.Config `koanf:"httpserver"`
	Postgres   postgres.Config   `koanf:"postgres"`
	Scheduler  scheduler.Config  `koanf:"scheduler"`
	GRPCServer grpcserver.Config `koanf:"grpcserver"`
}

func Load(path string) Config {
	var k = koanf.New(".")
	k.Load(file.Provider(path), yaml.Parser())

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}
	return cfg
}
