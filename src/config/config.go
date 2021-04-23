package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogo/protobuf/jsonpb"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"

	"github.com/will-rowe/herald/src/helpers"
	"github.com/will-rowe/herald/src/version"
)

var (

	// DefaultConfigName is the default config file name.
	DefaultConfigName = "herald-config"

	// DefaultConfigType of file for the config file.
	DefaultConfigType = "json"

	// DefaultConfigDir is the default location for the config file.
	DefaultConfigDir = getHome()

	// DefaultConfigPath for the config file.
	DefaultConfigPath = fmt.Sprintf("%s/%s.%s", DefaultConfigDir, DefaultConfigName, DefaultConfigType)

	// DefaultServerlog file path.
	DefaultServerlog = fmt.Sprintf("%s/herald-server.log", DefaultConfigDir)

	// DefaultManifestURL for the ARTIC primer schemes.
	DefaultManifestURL = "https://raw.githubusercontent.com/artic-network/primer-schemes/master/schemes_manifest.json"

	// ErrInvalidPath is used when the config file path is bad or doesn't exist.
	ErrInvalidPath = fmt.Errorf("invalid config filepath")

	// DefaultConfig is the basic info, filled on first run of Herald.
	DefaultConfig = &Config{
		Filepath:         DefaultConfigPath,
		Fileformat:       DefaultConfigType,
		User:             &User{},
		Version:          version.VERSION,
		Serverlog:        DefaultServerlog,
		ArticManifestURL: DefaultManifestURL,
	}
)

// Write will write a Config to disk.
func (config *Config) Write() error {
	if len(config.GetFilepath()) == 0 {
		return ErrInvalidPath
	}
	fh, err := os.Create(config.GetFilepath())
	defer fh.Close()
	d, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}
	_, err = fh.Write(d)
	return err
}

// GetJSONDump returns a string dump of the config in JSON
func (config *Config) GetJSONDump() string {

	// convert to JSON
	buf := &bytes.Buffer{}
	jsonMarshaller := jsonpb.Marshaler{
		EnumsAsInts:  false, // Whether to render enum values as integers, as opposed to string values.
		EmitDefaults: false, // Whether to render fields with zero values
		Indent:       "\t",  // A string to indent each level by
		OrigName:     false, // Whether to use the original (.proto) name for fields
	}
	jsonMarshaller.Marshal(buf, config)
	return string(buf.Bytes())
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
		if err := generateDefault(absConfigDir); err != nil {
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
func generateDefault(path string) error {

	// use a default directory if non provided
	if len(path) != 0 {
		DefaultConfig.Filepath = fmt.Sprintf("%s/%s.%s", path, DefaultConfigName, DefaultConfigType)
		DefaultConfig.Serverlog = fmt.Sprintf("%s/herald-server.log", path)
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
