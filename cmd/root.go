// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"io"
	"os"
	"strings"

	"gopkg.in/natefinch/lumberjack.v2"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configPath            = "$HOME"
	configFile            = ".dad"
	prefixForEnvVariables = "dad"
)

var cfgFile string

// RootCmd This represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "D.A.D",
	Short: "Deployment analytics dashboards",
	Long:  `D.A.D is a web application which manage the  deployment of services`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initLogger, initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.dad.yaml)")
	RootCmd.PersistentFlags().StringP("level", "l", "warning", "Choose the logger level: debug, info, warning, error, fatal, panic")
	RootCmd.PersistentFlags().Int("log-max-size", 100, "Max log file size in megabytes")
	RootCmd.PersistentFlags().Int("log-max-age", 30, "Max log file age in days")
	RootCmd.PersistentFlags().Int("log-max-backups", 3, "Max backup files to keep")
	_ = viper.BindPFlag("level", RootCmd.PersistentFlags().Lookup("level"))
	_ = viper.BindPFlag("log.max-size", RootCmd.PersistentFlags().Lookup("log-max-size"))
	_ = viper.BindPFlag("log.max-age", RootCmd.PersistentFlags().Lookup("log-max-age"))
	_ = viper.BindPFlag("log.max-backups", RootCmd.PersistentFlags().Lookup("log-max-backups"))
}

func initLogger() {
	output := io.MultiWriter(os.Stdout, &lumberjack.Logger{
		Filename:   "./logs/dad.log",
		MaxSize:    viper.GetInt("log.max-size"),
		MaxBackups: viper.GetInt("log.max-backups"),
		MaxAge:     viper.GetInt("log.max-age"),
	})
	log.SetOutput(output)

	level, err := log.ParseLevel(viper.GetString("level"))
	if err != nil {
		level = log.WarnLevel
		log.WithError(err).WithField("defaultLevel", level).Warn("Invalid log level, using default")
	}
	log.SetLevel(level)

	log.SetFormatter(&log.TextFormatter{})
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetEnvPrefix(prefixForEnvVariables)                        // Prefix for every env variables
	viper.SetConfigName(configFile)                                  // name of config file (without extension)
	viper.AddConfigPath(configPath)                                  // adding home directory as first search path
	viper.AutomaticEnv()                                             // read in environment variables that match
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_")) // Replace "." and "-" by "_" for env variable lookup

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err == nil {
		log.WithField("configFile", viper.ConfigFileUsed()).Info("Using provided config file")
	} else {
		log.WithError(err).Warn("Error with provided config file")
	}
}
