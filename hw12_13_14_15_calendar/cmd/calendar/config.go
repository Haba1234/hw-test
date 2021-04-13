package main

import (
	"context"

	"github.com/heetch/confita"
	"github.com/heetch/confita/backend/file"
)

type Config struct {
	Logger  LoggerConf
	Storage StorageConf
	Server  ServerConf
}

type LoggerConf struct {
	Level string `config:"level"`
	Path  string `config:"path"`
}

type StorageConf struct {
	Type     string `config:"type"`
	Host     string `config:"host"`
	PortDB   string `config:"portDB"`
	User     string `config:"user"`
	Password string `config:"password"`
	SSLMode  string `config:"SSLMode"`
}

type ServerConf struct {
	Address  string `config:"address"`
	HTTPPort string `config:"httpPort"`
	GRPCPort string `config:"grpcPort"`
}

func NewConfig(path string) (*Config, error) {
	// default values
	cfg := Config{
		Logger: LoggerConf{
			Level: "INFO",
			Path:  "./bin/logfile.log",
		},
		Storage: StorageConf{
			Type:     "memory",
			Host:     "localhost",
			PortDB:   "5432",
			User:     "postgres",
			Password: "",
			SSLMode:  "false",
		},
		Server: ServerConf{
			Address:  "localhost",
			HTTPPort: "8080",
			GRPCPort: "8090",
		},
	}

	loader := confita.NewLoader(
		file.NewBackend(path),
	)
	if err := loader.Load(context.Background(), &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
