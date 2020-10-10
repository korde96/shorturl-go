package config

import (
	"github.com/spf13/viper"
)

type (
	Config struct {
		Server   Server   `mapstructure:"server"`
		ASConfig ASConfig `mapstructure:"aerospike"`
	}

	Server struct {
		Host string `mapstructure:"host"`
		Port int    `mapstructure:"port"`
	}

	ASConfig struct {
		Host      string `mapstructure:"host"`
		Port      int    `mapstructure:"port"`
		Namespace string `mapstructure:"namespace"`
		TTL       int    `mapstructure:"ttl"`
	}
)

func LoadConfig() (Config, error) {
	var config Config
	vCfg := viper.GetViper()

	vCfg.AddConfigPath("./")
	vCfg.SetConfigName("config")

	err := vCfg.ReadInConfig()
	if err != nil {
		return config, err
	}

	err = vCfg.UnmarshalExact(&config)
	//add validation on the config
	return config, err
}
