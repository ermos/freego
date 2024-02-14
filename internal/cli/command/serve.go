package command

import (
	"fmt"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/proxy/controller"
	"github.com/ermos/freego/internal/proxy/logger"
	"github.com/ermos/freego/internal/proxy/service"
	"github.com/spf13/cobra"
	"net/http/httputil"
	"time"
)

type Serve struct {
	withDocker bool
}

func (d *Serve) Flags(cmd *cobra.Command) {
	cmd.Flags().BoolVar(
		&d.withDocker,
		"docker",
		false,
		"enable docker mode",
	)
}

func (d *Serve) Execute(cmd *cobra.Command, args []string) error {
	for {
		if err := config.Reload(); err != nil {
			logger.Error("serve/init", "config reload error", err)
			d.retry(20)
			continue
		}

		proxy := &httputil.ReverseProxy{
			Director:     controller.ProxyController{}.Handler(),
			ErrorHandler: controller.ErrorController{}.Handler(),
		}

		err := service.Start(service.ServiceContext{
			Proxy:      proxy,
			WithDocker: d.withDocker,
		})
		if err != nil {
			logger.Error("serve/init", "service start error", err)
			d.retry(20)
			continue
		}
	}
}

func (d *Serve) retry(sec int) {
	logger.Error("serve/init", fmt.Sprintf("retry in %d seconds...", sec))
	time.Sleep(time.Duration(sec) * time.Second)
}
