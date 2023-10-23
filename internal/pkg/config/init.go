package config

import (
	"github.com/spf13/viper"
	"os"
	"path/filepath"
)

func Init() (err error) {
	configDir, err := GetDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(configDir, "config.yaml")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(configDir, os.ModePerm)
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(configDir, "certificates"), os.ModePerm)
		if err != nil {
			return err
		}

		err = viper.WriteConfigAs(configPath)
		if err != nil {
			return err
		}
	}

	viper.AddConfigPath(filepath.Dir(configPath))
	viper.AutomaticEnv()

	return viper.ReadInConfig()
}

func Reload() (err error) {
	return viper.ReadInConfig()
}
