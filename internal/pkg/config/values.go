package config

import "github.com/spf13/viper"

func Set(key string, value any) {
	viper.Set(key, value)
}

func GetString(key string) string {
	return viper.GetString(key)
}

func Write() {
	viper.WriteConfig()
}
