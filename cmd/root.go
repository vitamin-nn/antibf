package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/vitamin-nn/antibf/internal/config"
	"github.com/vitamin-nn/antibf/internal/logger"
)

func Execute() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config file read error: %v", err)
	}

	err = logger.Init(cfg.Log)
	if err != nil {
		log.Fatalf("initialize logger error: %v", err)
	}

	rootCmd := &cobra.Command{
		Use:   "antibf",
		Short: "Antibruteforce grpc-service",
	}

	rootCmd.AddCommand(serverCmd(cfg))
	rootCmd.AddCommand(cliCmd(cfg))

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("execute cmd: %v", err)
	}
}
