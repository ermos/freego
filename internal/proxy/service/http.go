package service

import (
	"errors"
	"github.com/ermos/freego/internal/proxy/channel"
	"net/http"
)

type Http struct {
	server *http.Server
}

func (h *Http) Start(ctx ServiceContext) error {
	h.server = &http.Server{
		Addr:    ":80",
		Handler: ctx.Proxy,
	}

	go func() {
		err := h.server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			channel.InternalError.Send(err)
		}
	}()

	return nil
}

func (h *Http) Close(ServiceContext) error {
	if h.server == nil {
		return nil
	}
	return h.server.Close()
}
