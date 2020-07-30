package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"

	"github.com/vitamin-nn/otus_anti_bruteforce/internal/config"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/grpc"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/logger"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/ratelimit"
	rlMemRepo "github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/ratelimit/inmemory"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting"
	sMemRepo "github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/inmemory"
	sPsqlRepo "github.com/vitamin-nn/otus_anti_bruteforce/internal/repository/setting/psql"
	"github.com/vitamin-nn/otus_anti_bruteforce/internal/usecase"
)

func main() {
	var cfgFileName string
	flag.StringVar(&cfgFileName, "config", "./configs/local.toml", "config filepath")
	cfg, err := config.Load(cfgFileName)
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}

	err = logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("initialize logger error: %v", err)
	}
	log.WithFields(cfg.Fields()).Info("Starting antibruteforce service")

	loginLimit := ratelimit.NewRateLimit(rlMemRepo.NewSlidingWindow(), cfg.RateLimit.Login, cfg.RateLimit.Duration)
	passwdLimit := ratelimit.NewRateLimit(rlMemRepo.NewSlidingWindow(), cfg.RateLimit.Password, cfg.RateLimit.Duration)
	ipLimit := ratelimit.NewRateLimit(rlMemRepo.NewSlidingWindow(), cfg.RateLimit.IP, cfg.RateLimit.Duration)

	sMemRepo := sMemRepo.NewSettingRepo()
	sPsqlRepo := sPsqlRepo.NewSettingRepo(cfg.PSQL.DSN)
	ctx := context.Background()
	err = sPsqlRepo.Connect(ctx)
	if err != nil {
		log.Fatalf("psql connect error: %v", err)
	}
	defer sPsqlRepo.Close()

	s := setting.NewSettingRepo(sMemRepo, sPsqlRepo)
	rUseCase := usecase.NewRateLimitUseCase(loginLimit, passwdLimit, ipLimit, s)
	grpcServer := grpc.NewAntibruteforceServer(rUseCase)
	go grpcServer.Run(cfg.GrpcServer.Addr)

	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt)
	log.Infof("graceful shutdown: %v", <-interruptCh)

	log.Info("finished main program")
}
