package config

import (
	"flag"

	"github.com/caarlos0/env/v6"
)

type Config struct {
	GophermartAddr string `env:"RUN_ADDRESS" envDefault:"localhost:8070"`
	DBConnect      string `env:"DATABASE_URI"`
	AccrualAddr    string `env:"ACCRUAL_SYSTEM_ADDRESS" envDefault:"http://localhost:8080"`
}

func NewConfig() (*Config, error) {
	conf := &Config{}

	if err := env.Parse(conf); err != nil {
		return nil, err
	}

	regStringVar(&conf.GophermartAddr, "a", conf.GophermartAddr, "gophermart address")
	regStringVar(&conf.DBConnect, "d", conf.DBConnect, "database connection")
	regStringVar(&conf.AccrualAddr, "r", conf.AccrualAddr, "broker address")
	flag.Parse()

	return conf, nil
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}
