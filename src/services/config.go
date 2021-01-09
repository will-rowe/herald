package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/will-rowe/herald/src/helpers"
)

var (

	// DefaultConfigName is the default config file name.
	DefaultConfigName = "herald"

	// DefaultConfigType of file for the config file.
	DefaultConfigType = "json"

	// DefaultConfigDir is the default location for the config file.
	DefaultConfigDir = getHome()

	// DefaultConfigPath for the config file.
	DefaultConfigPath = fmt.Sprintf("%s/%s.%s", DefaultConfigDir, DefaultConfigName, DefaultConfigType)

	// ErrInvalidPath is used when the config file path is bad or doesn't exist.
	ErrInvalidPath = fmt.Errorf("invalid config filepath")

	// DefaultConfig is the basic info, filled on first run of Herald.
	DefaultConfig = &Config{
		Filepath:   DefaultConfigPath,
		Fileformat: DefaultConfigType,
		User:       &User{Name: "will rowe"},
	}
)

// Write will write a Config to disk.
func (x *Config) Write() error {
	if len(x.GetFilepath()) == 0 {
		return ErrInvalidPath
	}
	fh, err := os.Create(x.GetFilepath())
	defer fh.Close()
	d, err := json.MarshalIndent(x, "", "\t")
	if err != nil {
		return err
	}
	_, err = fh.Write(d)
	return err
}

// InitConfig reads in the config file
// or generates a new one if not found.
//
// It will return a pointer to a copy
// of the config in memory.
//
// NOTE: any updates to the returned
// config will need to be written back
// to the file on disk by the user.
func InitConfig(configDir string) (*Config, error) {

	// use the default config directory if non given
	if configDir == "" {
		configDir = DefaultConfigDir
	}
	absConfigDir, err := filepath.Abs(configDir)
	if err != nil {
		return nil, err
	}

	// put in a sanity check to strip whitespace from path (looking at you Google Drive)
	absConfigDir = strings.ReplaceAll(absConfigDir, " ", "\\ ")

	// if the config dir doesn't exist, add it now
	if err := helpers.CheckDirExists(absConfigDir); err != nil {
		if err := os.MkdirAll(absConfigDir, 0777); err != nil {
			return nil, err
		}
	}

	// if the config file doesn't exist, create a default one
	if exists := helpers.CheckFileExists(fmt.Sprintf("%s/%s.%s", absConfigDir, DefaultConfigName, DefaultConfigType)); !exists {
		if err := generateDefault(fmt.Sprintf("%s/%s.%s", absConfigDir, DefaultConfigName, DefaultConfigType)); err != nil {
			return nil, err
		}
	}

	// search for default config with viper
	viper.AddConfigPath(absConfigDir)
	viper.SetConfigName(DefaultConfigName)
	viper.SetConfigType(DefaultConfigType)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	// return a copy of the loaded config
	return dumpConfig()
}

// generateDefault will generate the default
// config on disk. If no filePath provided,
// it will use the DefaultConfigPath.
func generateDefault(filePath string) error {

	// use a default directory if non provided
	if len(filePath) != 0 {
		DefaultConfig.Filepath = filePath
	}

	// write the default config to disk
	return DefaultConfig.Write()
}

// resetDefault will remove any existing config and
// replace it with the default one.
//
// NOTE: the caller must reload
// the config into viper
func resetDefault(configPath string) error {

	// remove the existing config if it exists
	if helpers.CheckFileExists(configPath) {
		if err := os.Remove(configPath); err != nil {
			return err
		}
	}

	// now generate the default and write it to disk
	return generateDefault(configPath)
}

// dumpConfig will unmarshall the config from
// Viper to a struct in memory.
func dumpConfig() (*Config, error) {
	c := &Config{}
	err := viper.UnmarshalExact(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// getHome is used to find the DefaultConfigDir.
func getHome() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return homeDir
}
