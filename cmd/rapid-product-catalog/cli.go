package main

import (
	"fmt"

	"github.com/roppenlabs/rapid-product-catalog/internal/config"

	"github.com/spf13/cobra"

	logger "github.com/roppenlabs/rapido-logger-go"
)

func initCLI() *cobra.Command {
	var cliCmd = &cobra.Command{
		Use:   "rapid-product-catalog",
		Short: "rapid-product-catalog CLI to manage the service",
	}

	cliCmd.AddCommand(startCommand())
	return cliCmd
}

func startCommand() *cobra.Command {
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts the service",
		Run: func(cmd *cobra.Command, args []string) {
			configConfig, err := config.NewConfig()
			if err != nil {
				panic(fmt.Errorf("failed to initialize config: %w", err))
			}
			config.SetConfig(configConfig)
			logger.Init(configConfig.Get().Log.Level)

			serverDependencies, err := InitDependencies()
			if err != nil {
				panic(fmt.Errorf("failed to initialize dependencies: %w", err))
			}
			serverDependencies.server.Run(serverDependencies.handlers)
		},
	}

	return startCmd
}
