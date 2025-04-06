// internal/config/config.go
package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	ClientID     string `mapstructure:"CLIENT_ID"`
	ClientSecret string `mapstructure:"CLIENT_SECRET"`
	RedirectURL  string `mapstructure:"REDIRECT_URL"`
}

func LoadConfig() (*Config, error) {
	// Set the file name of the configurations file
	viper.SetConfigFile(".envrc")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()

	// Read in the config file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("fatal error config file: %s", err)
	}

	// Unmarshal the config into the Config struct
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &cfg, nil
}
