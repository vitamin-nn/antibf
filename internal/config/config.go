package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func Load(filepath string) (c *Config, err error) {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("toml")
	setDefaults()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err = viper.Unmarshal(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

func setDefaults() {
	viper.SetDefault("log.log_level", "info")
}

type Config struct {
	Log        Log
	GrpcServer GrpcServer `mapstructure:"grpc_server"`
	PSQL       PSQL
	RateLimit  RateLimit `mapstructure:"ratelimits"`
}

type PSQL struct {
	DSN string
}

type Log struct {
	LogLevel string `mapstructure:"log_level"`
}

type GrpcServer struct {
	Addr string
}

type RateLimit struct {
	Login    int
	Password int
	IP       int
	Duration int64
}

func (cfg *Config) Fields() log.Fields {
	return log.Fields{
		"server_addr":   cfg.GrpcServer.Addr,
		"login_rate":    cfg.RateLimit.Login,
		"passwd_rate":   cfg.RateLimit.Password,
		"ip_rate":       cfg.RateLimit.IP,
		"rate_duration": cfg.RateLimit.Duration,
		"log_level":     cfg.Log.LogLevel,
	}
}
