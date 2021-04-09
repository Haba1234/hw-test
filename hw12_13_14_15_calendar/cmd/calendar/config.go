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
	SSLMode  bool   `config:"SSLMode"`
}

type ServerConf struct {
	Address string `config:"address"`
	Port    string `config:"port"`
}

func NewConfig(path string) (*Config, error) {
	cfg := &Config{}
	loader := confita.NewLoader(
		file.NewBackend(path),
	)
	if err := loader.Load(context.Background(), cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
