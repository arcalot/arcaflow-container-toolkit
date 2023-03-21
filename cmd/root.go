/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.arcalot.io/log"
)

var cfgFile string
var verbosity bool
var rootLogger log.Logger

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "arcaflow-plugin-image-builder",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		log_level := log.LevelInfo
		if verbosity {
			log_level = log.LevelDebug
		}
		ConfigureLogger(&rootLogger, log_level, log.DestinationStdout, os.Stdout)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		rootLogger.Errorf("root command failed", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.act.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().BoolVarP(&verbosity, "verbosity", "v", false, "verbose debugging log messages")
	ConfigureLogger(&rootLogger, log.LevelInfo, log.DestinationStdout, os.Stdout)
}

func ConfigureLogger(logger *log.Logger, level log.Level, dest log.Destination, w io.Writer) {
	logConfig := log.Config{
		Level:       level,
		Destination: dest,
		Stdout:      w,
	}
	*logger = log.New(logConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".act")
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		rootLogger.Infof("Using config file:%s\n", viper.ConfigFileUsed())
	} else {
		rootLogger.Errorf("Did not find .act config file")
		os.Exit(1)
	}

	viper.AutomaticEnv() // read in environment variables that match
}
