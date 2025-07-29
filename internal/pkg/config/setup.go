package config

import "github.com/spf13/viper"

func Setup() {
	viper.SetConfigName("tunes_config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.WriteConfig()
		}
	}
}
