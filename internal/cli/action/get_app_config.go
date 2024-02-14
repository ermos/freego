package action

import (
	"encoding/base64"
	"fmt"
	"github.com/ermos/freego/internal/cli/model"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"strings"
)

func GetAppConfig(customConfigFile string) (model.AppConfig, error) {
	var c model.AppConfig

	dir, err := os.Getwd()
	if err != nil {
		return c, err
	}

	cfgPath := filepath.Join(dir, "freego.yaml")
	if customConfigFile != "" {
		cfgPath = filepath.Join(dir, customConfigFile)
	}

	configData, err := os.ReadFile(cfgPath)
	if err != nil {
		if customConfigFile == "" {
			cfgPath = filepath.Join(dir, "freego.yml")
			configData, err = os.ReadFile(cfgPath)
		}
		if err != nil {
			return c, err
		}
	}

	err = yaml.Unmarshal(configData, &c)
	if err != nil {
		return c, err
	}

	for domain, content := range c.Domains {
		if content.Host == "" {
			content.Host = "127.0.0.1"
		}

		if content.Port == 0 {
			return c, fmt.Errorf("%s: port is required", domain)
		}

		if content.Port == 80 {
			return c, fmt.Errorf("%s: port 80 is reserved", domain)
		}

		c.Domains[domain] = content
	}

	c.Link = strings.ReplaceAll(
		strings.ToLower(
			base64.StdEncoding.EncodeToString([]byte(cfgPath)),
		),
		"=",
		"",
	)

	return c, err
}
