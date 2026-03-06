package config

import (
	"fmt"

	"github.com/Hans-Kerman/go-book-lending/backend/types"
	"github.com/spf13/viper"
)

var AppConfig *types.Config

func InitConfig() error {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("viper reading config file: %w", err)
	}

	if err = viper.Unmarshal(AppConfig); err != nil {
		return fmt.Errorf("viper unmarshal config file: %w", err)
	}

	return nil
}
