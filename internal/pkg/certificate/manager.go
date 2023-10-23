package certificate

import (
	"crypto/tls"
	"github.com/ermos/freego/internal/pkg/config"
	"os"
	"path/filepath"
	"strings"
)

var certificates = make(map[string]*tls.Certificate)

func Load() (err error) {
	configDir, err := config.GetDir()
	if err != nil {
		return err
	}

	err = filepath.Walk(filepath.Join(configDir, "certificates"),
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !strings.HasSuffix(path, "cert.pem") {
				return nil
			}

			cert, err := tls.LoadX509KeyPair(path, strings.Replace(path, "cert.pem", "key.pem", 1))
			if err != nil {
				return err
			}

			domain := filepath.Base(strings.Replace(path, "cert.pem", "", 1))

			certificates[domain] = &cert

			return nil
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func GetCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	return certificates[info.ServerName], nil
}

func GetAll() []tls.Certificate {
	var cs []tls.Certificate

	for _, c := range certificates {
		cs = append(cs, *c)
	}

	return cs
}
