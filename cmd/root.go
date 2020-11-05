// Copyright © 2020 github.com/ptrkrlsrd
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
	"fmt"
	"os"

	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

var (
	cfgFile string
	service acache.Server
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acache",
	Short: "Simple API cacher and server",
}

// HandleError Handle and error by printing the error and returning Exit code 1
func HandleError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		HandleError(err)
	}
}

func init() {
	// Set the flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.config/acache/acache.json", "Config file")
	rootCmd.PersistentFlags().StringP("database", "d", "~/.config/acache/acache.db", "Database")

	// Initialize the database and config
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initDB)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".config/acache" (without extension).
		configPath, err := configPath()
		if err != nil {
			HandleError(err)
		}

		viper.AddConfigPath(configPath)
		viper.SetConfigName("config")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func initDB() {
	path := rootCmd.Flag("database").Value.String()
	expandedPath, err := tilde.Expand(path)
	db, err := acache.NewDB(expandedPath)
	if err != nil {
		HandleError(err)
	}

	storage, err := acache.NewStorage("acache", expandedPath, db)
	if err != nil {
		HandleError(err)
	}

	service = acache.Server{Storage: storage}
	service.Storage.LoadRoutes()
}

func configPath() (string, error) {
	path := rootCmd.Flag("config").Value.String()
	return tilde.Expand(path)
}
