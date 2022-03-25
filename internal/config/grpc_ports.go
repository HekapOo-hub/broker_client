package config

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Ports struct {
	Position string `env:"POSITION_PORT" envDefault:":50005"`
	Order    string `env:"ORDER_PORT" envDefault:":50004"`
}

func GetPorts() (*Ports, error) {
	cfg := Ports{}
	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("error with parsing env variables in position port config %w", err)
	}
	return &cfg, nil
}
