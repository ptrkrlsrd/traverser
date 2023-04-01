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
	cfgFile        string
	databasePath   string
	yamlFilePath   string
	useYamlStorage bool
	server         acache.Server
	storage        acache.RouteStorer
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "~/.config/acache/acache.json", "Config file")
	rootCmd.PersistentFlags().StringVar(&databasePath, "d", "~/.config/acache/", "Database")
	rootCmd.PersistentFlags().BoolVarP(&useYamlStorage, "use-yaml", "y", false, "Use YAML storage")
	rootCmd.PersistentFlags().StringVar(&yamlFilePath, "yaml-path", "./routes.yaml", "Use YAML storage")
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initStorage)
	cobra.OnInitialize(initServer)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".config/acache" (without extension).
		configPath, err := tilde.Expand(cfgFile)
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

func initStorage() {
	if useYamlStorage {
		initFileStore()
	} else {
		initBadgerDB()
	}
}

func initBadgerDB() {
	expandedConfigPath, err := tilde.Expand(databasePath)
	HandleError(err)

	checkOrCreateFolder(expandedConfigPath)
	HandleError(err)

	db, err := acache.NewBadgerDB(path.Join(expandedConfigPath, "acache.db"))
	HandleError(err)

	storage, err = acache.NewBadgerStorage(db)
	HandleError(err)
}

func initServer() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	server = acache.NewServer(storage, router)
}

func initFileStore() {
	var err error
	storage, err = acache.NewYAMLStorage(yamlFilePath)
	if err != nil {
		HandleError(err)
	}
}

func checkOrCreateFolder(expandedConfigPath string) error {
	if _, err := os.Stat(expandedConfigPath); os.IsNotExist(err) {
		return os.Mkdir(expandedConfigPath, os.ModePerm)
	}
	return nil
}
