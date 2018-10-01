package main

import (
	"github.com/awslabs/aws-service-operator/pkg/logger"
	"github.com/awslabs/aws-service-operator/pkg/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
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

		signalChan := make(chan os.Signal, 1)
		stopChan := make(chan struct{})
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		server.New(config).Run(stopChan)

		for {
			select {
			case <-signalChan:
				logger.Info("shutdown signal received, exiting...")
				close(stopChan)
				return
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
