package config

import (
	"github.com/spf13/viper"
	"os"
	"os/user"
	"path/filepath"
)

func Init() (err error) {
	var configDir string

	if os.Getenv("SUDO_USER") != "" {
		var u *user.User

		u, err = user.Lookup(os.Getenv("SUDO_USER"))
		if err != nil {
			return err
		}

		configDir = filepath.Join(u.HomeDir, ".config")
	} else {
		configDir, err = os.UserConfigDir()
		if err != nil {
			return err
		}
	}

	configPath := filepath.Join(configDir, "freego/config.yaml")

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
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
