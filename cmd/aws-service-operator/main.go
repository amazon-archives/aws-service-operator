package main

import (
	"fmt"
	homedir "github.com/mitchellh/go-homedir"

	"github.com/awslabs/aws-service-operator/pkg/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var (
	// cfgFile, kubeConfig, awsRegion all help support passed in flags into the server
	cfgFile, kubeconfig, awsRegion, logLevel, logFile, resources, clusterName, bucket, accountID string

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "aws-operator",
		Short: "AWS Operator manages your AWS Infrastructure using CRDs and Operators",
		Long: `AWS Operator manages your AWS Infrastructure using CRDs and Operators. 
With a single manifest file you can now model both the application and the resource necessary to run it.`,
		Run: func(c *cobra.Command, _ []string) {
			c.Help()
		},
	}
)

func main() {
	cobra.OnInitialize(initConfig)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "f", "Config file (default is $HOME/.aws-operator.yaml)")
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to local kubeconfig file (mainly used for development)")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "us-west-2", "AWS Region for resources to be created in")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "Info", "Log level for the CLI")
	rootCmd.PersistentFlags().StringVarP(&logFile, "logfile", "", "", "Log level for the CLI")
	rootCmd.PersistentFlags().StringVarP(&resources, "resources", "", "s3bucket,dynamodb", "Comma delimited list of CRDs to deploy")
	rootCmd.PersistentFlags().StringVarP(&clusterName, "cluster-name", "i", "aws-operator", "Cluster name for the Application to run as, used to label the Cloudformation templated to avoid conflict")
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "aws-operator", "To configure the operator you need a base bucket to contain the resources")
	rootCmd.PersistentFlags().StringVarP(&accountID, "account-id", "a", "", "AWS Account ID, this is used to configure outputs and operate on the proper account.")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("resources", rootCmd.PersistentFlags().Lookup("resources"))
	viper.BindPFlag("clustername", rootCmd.PersistentFlags().Lookup("cluster-name"))
	viper.BindPFlag("bucket", rootCmd.PersistentFlags().Lookup("bucket"))
	viper.BindPFlag("accountid", rootCmd.PersistentFlags().Lookup("account-id"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".aws-operator")
	}

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getConfig() (*config.Config, error) {
	resourcesList := strings.Split(resources, ",")
	config := &config.Config{
		Region:     awsRegion,
		Kubeconfig: kubeconfig,
		LoggingConfig: &config.LoggingConfig{
			File:              logFile,
			Level:             logLevel,
			FullTimestamps:    true,
			DisableTimestamps: false,
		},
		Resources:   resourcesList,
		ClusterName: clusterName,
		Bucket:      bucket,
		AccountID:   accountID,
	}

	return config, nil
}
