package config

import (
	"os"
	"os/user"
	"path/filepath"
)

func GetDir() (string, error) {
	if os.Getenv("SUDO_USER") != "" {
		u, err := user.Lookup(os.Getenv("SUDO_USER"))
		if err != nil {
			return "", err
		}

		return filepath.Join(u.HomeDir, ".config/freego"), nil
	}

	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configDir, "freego"), nil
}
