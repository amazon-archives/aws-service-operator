package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/awslabs/aws-service-operator/pkg/logger"
	"github.com/awslabs/aws-service-operator/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the operator",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		config, err := getConfig()
		if err != nil {
			logrus.Fatalf("%s", err)
		}

		logger, err := logger.Configure(config.LoggingConfig)
		if err != nil {
			logrus.Fatalf("Failed to configure logging: '%s'" + err.Error())
		}
		config.Logger = logger

		ctx, cancel := context.WithCancel(context.Background())
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		go server.New(config).Run(ctx)

		<-signalChan
		logger.Info("shutdown signal received, exiting...")
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
