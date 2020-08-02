package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "nagiosapi",
		Short: "Provides an HTTP API to Nagios",
	}
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigType(cfgFile)
	} else {
		viper.SetConfigName("nagios-api")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("/etc/nagios-api/")
		viper.AddConfigPath(".")
	}

	viper.SetEnvPrefix("NAGAPI")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("unable to read config: %v", err)
		}
	}
}
