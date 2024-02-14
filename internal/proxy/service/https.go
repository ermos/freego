package service

import (
	"crypto/tls"
	"errors"
	"github.com/ermos/freego/internal/pkg/certificate"
	"github.com/ermos/freego/internal/proxy/channel"
	"net/http"
)

type Https struct {
	server *http.Server
}

func (h *Https) Start(ctx ServiceContext) error {
	err := certificate.Load()
	if err != nil {
		return err
	}

	h.server = &http.Server{
		Addr:    ":443",
		Handler: ctx.Proxy,
		TLSConfig: &tls.Config{
			Certificates:   certificate.GetAll(),
			GetCertificate: certificate.GetCertificate,
		},
	}

	go func() {
		err = h.server.ListenAndServeTLS("", "")
		if !errors.Is(err, http.ErrServerClosed) {
			channel.InternalError.Send(err)
		}
	}()

	return nil
}

func (h *Https) Close(ServiceContext) error {
	if h.server == nil {
		return nil
	}
	return h.server.Close()
}
