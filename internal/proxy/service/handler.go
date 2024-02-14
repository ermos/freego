package service

import (
	"github.com/ermos/freego/internal/proxy/channel"
	"github.com/ermos/freego/internal/proxy/logger"
	"net/http/httputil"
)

type Service interface {
	Start(ctx ServiceContext) error
	Close(ctx ServiceContext) error
}

type ServiceContext struct {
	Proxy      *httputil.ReverseProxy
	WithDocker bool
}

var handler = []Service{
	&Http{},
	&Https{},
	&Config{},
	&Docker{},
}

func Start(ctx ServiceContext) error {
	if err := startAll(ctx); err != nil {
		return err
	}

	select {
	case err := <-channel.InternalError.Listen():
		return err
	case <-channel.Reload.Listen():
		logger.Info("serve/service", "config updated ! Reloading...")
		if err := closeAll(ctx); err != nil {
			return err
		}
	}

	return nil
}

func startAll(ctx ServiceContext) error {
	for _, h := range handler {
		if err := h.Start(ctx); err != nil {
			return err
		}
	}

	return nil
}

func closeAll(ctx ServiceContext) error {
	for _, h := range handler {
		if err := h.Close(ctx); err != nil {
			return err
		}
	}

	return nil
}
