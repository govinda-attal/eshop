package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-logrusutil/logrusutil/logctx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	cfg     Config

	log = logctx.Default
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "eshop",
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/app-cfg.yaml", "application config file")

	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	viper.SetConfigFile(cfgFile)
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal("unable to unmarshal configuration locally")
	}
	lvl, _ := logrus.ParseLevel(cfg.Log.Level)
	log.Logger.SetLevel(lvl)
	log.Logger.SetFormatter(&logrus.JSONFormatter{})
	log.WithField("config", cfg).Info("application configuration")
}
