package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	configDirPerm = 0750
)

var root = &cobra.Command{
	Use: "wms",
	Long: `This program helps you to generate images via web map services.

Configuration file: $HOME/wms-config/.wms.yaml
Further informations and examples: https://github.com/wroge/wms`,
	SilenceUsage: true,
}

var version = "No Version Provided"

var versionCommand = &cobra.Command{
	Use:   "version",
	Short: "Show Version",
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println(version)
	},
}

func init() {
	cobra.OnInitialize(initConfig)
	root.PersistentFlags().Bool("help", false, "Help about any command")
	root.AddCommand(versionCommand)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	viper.AddConfigPath(filepath.Join(home, "wms-config"))
	viper.SetConfigName(".wms")
	if err := viper.ReadInConfig(); err != nil {
		createConfigIfNotExists(home)
	}
}

// createConfigIfNotExists creates the config directory and file if they don't exist.
func createConfigIfNotExists(home string) {
	configDir := filepath.Join(home, "wms-config")
	configFile := filepath.Join(configDir, ".wms.yaml")
	
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, configDirPerm); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}
		// #nosec G304 - This is the standard config path for the application
		if _, err := os.Create(configFile); err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println("Can't read config file")
		os.Exit(1)
	}
}

// Execute root command.
func Execute(v string) {
	version = v
	err := root.Execute()
	if err != nil {
		os.Exit(1)
	}
}
