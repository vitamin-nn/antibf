package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

func Load() (*Config, error) {
	cfg := new(Config)
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

type Config struct {
	Log        Log
	GrpcServer GrpcServer
	PSQL       PSQL
	RateLimit  RateLimit
	WBSetting  WhiteBlackSetting
}

type PSQL struct {
	User     string `env:"POSTGRES_USER,required"`
	Password string `env:"POSTGRES_PASSWORD,required"`
	DB       string `env:"POSTGRES_DB,required"`
	DBHost   string `env:"POSTGRES_DB_HOST,required"`
	Port     int    `env:"POSTGRES_PORT,required"`
}

type Log struct {
	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

type GrpcServer struct {
	Addr string `env:"GRPC_SERVER_ADDR,required"`
}

type RateLimit struct {
	Login    int           `env:"RATELIMIT_LOGIN,required"`
	Password int           `env:"RATELIMIT_PASSWORD,required"`
	IP       int           `env:"RATELIMIT_IP,required"`
	Duration time.Duration `env:"RATELIMIT_DURATION,required"`
}

type WhiteBlackSetting struct {
	CacheUpdInterval time.Duration `env:"SETTING_RELOAD_INTERVAL,required"`
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

func (cpg *PSQL) GetDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cpg.DBHost, cpg.Port, cpg.User, cpg.Password, cpg.DB)
}
