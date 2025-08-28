package config

import "github.com/spf13/viper"

func SetYtDlpPath(value any) {
	viper.Set(ytDlpPath, value)
}

func SetYtDlpRelease(value any) {
	viper.Set(ytDlpRelease, value)
}

func GetYtDlpPath() string {
	return viper.GetString(ytDlpPath)
}

func GetYtDlpRelease() string {
	return viper.GetString(ytDlpRelease)
}

func Write() {
	viper.WriteConfig()
}
