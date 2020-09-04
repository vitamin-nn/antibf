package cmd

import (
	"context"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vitamin-nn/antibf/internal/config"
	"github.com/vitamin-nn/antibf/internal/grpc"
	"github.com/vitamin-nn/antibf/internal/ratelimit"
	"github.com/vitamin-nn/antibf/internal/repository/setting"
	sPsqlRepo "github.com/vitamin-nn/antibf/internal/repository/setting/psql"
	"github.com/vitamin-nn/antibf/internal/usecase"
)

func serverCmd(cfg *config.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Starts antibruteforce grpc-server",
		Run: func(cmd *cobra.Command, args []string) {
			log.WithFields(cfg.Fields()).Info("Starting antibruteforce service")

			ctx, cancel := context.WithCancel(context.Background())
			loginLimit := ratelimit.NewRateLimit(ctx, cfg.RateLimit.Login, cfg.RateLimit.Duration)
			passwdLimit := ratelimit.NewRateLimit(ctx, cfg.RateLimit.Password, cfg.RateLimit.Duration)
			ipLimit := ratelimit.NewRateLimit(ctx, cfg.RateLimit.IP, cfg.RateLimit.Duration)

			sPsqlRepo := sPsqlRepo.NewSettingRepo(cfg.PSQL.GetDSN())
			err := sPsqlRepo.Connect(ctx)
			if err != nil {
				log.Fatalf("psql connect error: %v", err)
			}
			defer sPsqlRepo.Close()

			s := setting.NewSettingRepo(ctx, sPsqlRepo, cfg.WBSetting.CacheUpdInterval)
			rUseCase := usecase.NewRateLimitUseCase(loginLimit, passwdLimit, ipLimit, s)
			sUseCase := usecase.NewSettingUseCase(s)
			grpcServer := grpc.NewAntibruteforceServer(rUseCase, sUseCase)

			err = grpcServer.Run(cfg.GrpcServer.Addr)
			if err != nil {
				log.Fatalf("grpc server running error: %v", err)
			}

			go func() {
				interruptCh := make(chan os.Signal, 1)
				signal.Notify(interruptCh, os.Interrupt)
				log.Infof("graceful shutdown: %v", <-interruptCh)
				cancel()
				log.Info("finished main program")
			}()
		},
	}
}
