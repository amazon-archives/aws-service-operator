package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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

		ctx, cancel := context.WithCancel(context.Background())
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		go server.New(config).Run(ctx)

		<-signalChan
		config.Logger.Info("shutdown signal received, exiting...")
		cancel()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
