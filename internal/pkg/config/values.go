package config

import "github.com/spf13/viper"

func SetYtDlpPath(value string) {
	viper.Set(ytDlpPath, value)
}

func SetYtDlpRelease(value string) {
	viper.Set(ytDlpRelease, value)
}

func GetYtDlpPath() string {
	return viper.GetString(ytDlpPath)
}

func GetYtDlpRelease() string {
	return viper.GetString(ytDlpRelease)
}

func GetMaxThreads() int {
	return viper.GetInt(maxThreads)
}

func SetMaxThreads(threads int) {
	viper.Set(maxThreads, threads)
}

func Write() {
	viper.WriteConfig()
}
