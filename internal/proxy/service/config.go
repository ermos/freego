package service

import (
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/proxy/channel"
	"github.com/fsnotify/fsnotify"
	"log"
)

type Config struct {
	watcher *fsnotify.Watcher
}

func (c *Config) Start(ServiceContext) (err error) {
	c.watcher, err = fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case event, ok := <-c.watcher.Events:
				if ok && event.Has(fsnotify.Write) {
					channel.Reload.Send()
					return
				}
			case watcherErr, ok := <-c.watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", watcherErr)
			}
		}
	}()

	return c.watcher.Add(config.GetConfigPath())
}

func (c *Config) Close(ServiceContext) error {
	return c.watcher.Close()
}
