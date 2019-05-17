package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/awslabs/aws-service-operator/pkg/config"
	"github.com/awslabs/aws-service-operator/pkg/logger"
	"github.com/awslabs/aws-service-operator/pkg/queue"
	goVersion "github.com/christopherhein/go-version"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
)

var (
	// cfgFile, masterURL, kubeConfig, awsRegion all help support passed in flags into the server
	cfgFile, masterURL, kubeconfig, awsRegion, logLevel, logFile, resources, clusterName, bucket, accountID, k8sNamespace string
	defaultNamespace                                                                                                      string

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
	rootCmd.PersistentFlags().StringVarP(&masterURL, "master-url", "u", "", "The address of the Kubernetes API server. Overrides any value in kubeconfig.")
	rootCmd.PersistentFlags().StringVarP(&kubeconfig, "kubeconfig", "k", "", "Path to local kubeconfig file (mainly used for development)")
	rootCmd.PersistentFlags().StringVarP(&awsRegion, "region", "r", "us-west-2", "AWS Region for resources to be created in")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "loglevel", "l", "Info", "Log level for the CLI")
	rootCmd.PersistentFlags().StringVarP(&logFile, "logfile", "", "", "Log file for the CLI")
	rootCmd.PersistentFlags().StringVarP(&resources, "resources", "", "cloudformationtemplate,dynamodb,ecrrepository,elasticache,s3bucket,snssubscription,snstopic,sqsqueue", "Comma delimited list of CRDs to deploy")
	rootCmd.PersistentFlags().StringVarP(&clusterName, "cluster-name", "i", "aws-operator", "Cluster name for the Application to run as, used to label the Cloudformation templated to avoid conflict")
	rootCmd.PersistentFlags().StringVarP(&bucket, "bucket", "b", "aws-operator", "To configure the operator you need a base bucket to contain the resources")
	rootCmd.PersistentFlags().StringVarP(&accountID, "account-id", "a", "", "AWS Account ID, this is used to configure outputs and operate on the proper account.")
	rootCmd.PersistentFlags().StringVarP(&k8sNamespace, "k8s-namespace", "", "", "Namespace to scope k8s API queries to. If left blank will default to all namespaces")
	rootCmd.PersistentFlags().StringVarP(&defaultNamespace, "default-namespace", "", "default", "The default namespace in which to look for CloudFormation templates")

	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("masterurl", rootCmd.PersistentFlags().Lookup("master-url"))
	viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	viper.BindPFlag("loglevel", rootCmd.PersistentFlags().Lookup("loglevel"))
	viper.BindPFlag("logfile", rootCmd.PersistentFlags().Lookup("logfile"))
	viper.BindPFlag("resources", rootCmd.PersistentFlags().Lookup("resources"))
	viper.BindPFlag("clustername", rootCmd.PersistentFlags().Lookup("cluster-name"))
	viper.BindPFlag("bucket", rootCmd.PersistentFlags().Lookup("bucket"))
	viper.BindPFlag("accountid", rootCmd.PersistentFlags().Lookup("account-id"))
	viper.BindPFlag("defaultnamespace", rootCmd.PersistentFlags().Lookup("default-namespace"))
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

func getConfig() (c config.Config, err error) {
	loggingConfig := config.LoggingConfig{
		File:              logFile,
		Level:             logLevel,
		FullTimestamps:    true,
		DisableTimestamps: false,
	}

	resourcesMap := map[string]bool{}
	for _, r := range strings.Split(resources, ",") {
		resourcesMap[r] = true
	}

	ec2Session, err := session.NewSession()
	metadata := ec2metadata.New(ec2Session)
	if awsRegion == "" {
		awsRegion, err = metadata.Region()
		if err != nil {
			return c, err
		}
	}

	sess, err := session.NewSession(&aws.Config{Region: aws.String(awsRegion)})
	if err != nil {
		return c, err
	}

	awsclientset, kubeclientset, restconfig, err := config.CreateContext(masterURL, kubeconfig)
	if err != nil {
		return c, err
	}

	logger, err := logger.Configure(loggingConfig)
	if err != nil {
		return c, err
	}

	queueURL, queueARN, err := queue.RegisterQueue(sess, clusterName, "cloudformation")
	if err != nil {
		return c, err
	}

	c = config.Config{
		Region:           awsRegion,
		Kubeconfig:       kubeconfig,
		MasterURL:        masterURL,
		Logger:           logger,
		Version:          goVersion.New(version, commit, date),
		AWSSession:       sess,
		LoggingConfig:    loggingConfig,
		AWSClientset:     awsclientset,
		KubeClientset:    kubeclientset,
		RESTConfig:       restconfig,
		Recorder:         config.CreateRecorder(logger, kubeclientset),
		Resources:        resourcesMap,
		ClusterName:      clusterName,
		Bucket:           bucket,
		AccountID:        accountID,
		DefaultNamespace: defaultNamespace,
		K8sNamespace:     k8sNamespace,
		QueueURL:         queueURL,
		QueueARN:         queueARN,
	}
	return c, nil
}
