package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/ptrkrlsrd/acache/pkg/acache"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	tilde "gopkg.in/mattes/go-expand-tilde.v1"
)

var (
	cfgFile string
	server  acache.Server
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "acache",
	Short: "API response recorder",
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
	err := rootCmd.Execute()
	HandleError(err)
}

func init() {
	// Set the flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.config/acache/acache.json", "Config file")
	rootCmd.PersistentFlags().StringP("database", "d", "~/.config/acache/", "Database")

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
		configPath, err := tilde.Expand(rootCmd.Flag("config").Value.String())
		HandleError(err)

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
	pathVariable := rootCmd.Flag("database").Value.String()
	expandedConfigPath, err := tilde.Expand(pathVariable)
	HandleError(err)

	checkOrCreateFolder(expandedConfigPath)
	HandleError(err)

	db, err := acache.NewBadgerDB(path.Join(expandedConfigPath, "acache.db"))
	HandleError(err)

	storage, err := acache.NewBadgerStorage(db)
	HandleError(err)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	server = acache.NewServer(storage, router)
}

func checkOrCreateFolder(expandedConfigPath string) error {
	if _, err := os.Stat(expandedConfigPath); os.IsNotExist(err) {
		return os.Mkdir(expandedConfigPath, os.ModePerm)
	}
	return nil
}
