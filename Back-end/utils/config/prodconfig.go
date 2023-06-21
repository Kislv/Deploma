package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ProdConfig struct {
	Database string `mapstructure:"database"`
	Feed     uint64 `mapstructure:"feed limit"`
	Genres   uint64 `mapstructure:"genres limit"`
}

var ProdConfigStore ProdConfig

const (
	prodFilename = "prodconfig.json"
	prodExt      = "json"
)

func (cfg *ProdConfig) FromJson() error {
	viper.AddConfigPath(configpath)
	viper.SetConfigName(prodFilename)
	viper.SetConfigType(prodExt)

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error in config reading")
		cfg.clear()
		return err
	}

	if err := viper.Unmarshal(&ProdConfigStore); err != nil {
		fmt.Println("error in config reading")
		cfg.clear()
		return err
	}

	return nil
}

func (cfg *ProdConfig) clear() {
	cfg.Database = ""
	cfg.Feed = 12
	cfg.Genres = 10
}
