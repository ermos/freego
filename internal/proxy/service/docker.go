package service

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/ermos/freego/internal/cli/action"
	"github.com/ermos/freego/internal/pkg/config"
	"github.com/ermos/freego/internal/proxy/logger"
	"github.com/moby/moby/client"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
)

type Docker struct {
	client *client.Client
}

func (d *Docker) Start(serviceCtx ServiceContext) (err error) {
	if !serviceCtx.WithDocker {
		return
	}

	d.client, err = client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return
	}

	ctx := context.Background()

	go func() {
		eventsChan, errs := d.client.Events(ctx, types.EventsOptions{})

		for {
			select {
			case event := <-eventsChan:
				if event.Type != events.ContainerEventType {
					continue
				}

				var container types.ContainerJSON
				domains := make(map[string]*config.ActiveDomain)

				container, err = d.client.ContainerInspect(ctx, event.Actor.ID)
				if err != nil {
					logger.Info(
						"serve/docker",
						fmt.Sprintf("remove hosts from %s (%s)", event.Actor.ID, event.Action),
					)

					err = action.RemoveHostsFromLink(event.Actor.ID)
					if err != nil {
						logger.Error(
							"serve/docker",
							fmt.Sprintf("error removing hosts from %s (%s)", event.Actor.ID, event.Action),
							err,
						)
					}

					return
				}

				if event.Action != "create" && event.Action != "update" {
					continue
				}

				for label, value := range container.Config.Labels {
					if strings.HasPrefix(label, "freego.") {
						part := strings.Split(label, ".")
						if len(part) != 3 {
							continue
						}

						if domains[part[1]] == nil {
							domains[part[1]] = &config.ActiveDomain{}
						}

						switch part[2] {
						case "hostname":
							domains[part[1]].Domain = value
						case "address":
							domains[part[1]].Host = value
						case "port":
							var port int

							port, err = strconv.Atoi(value)
							if err == nil {
								domains[part[1]].Port = port
							}
						}
					}
				}

				// @TODO if port not renseigned, use container port
				for _, domain := range domains {
					domain.Link = event.Actor.ID
				}

				fmt.Println(domains)

				logger.Info(
					"serve/docker",
					fmt.Sprintf("container %s %s", event.Actor.ID, event.Action),
				)

			case err = <-errs:
				if err != nil {
					log.Printf("Error receiving Docker events: %s\n", err)
				}
			}
		}
	}()

	return
}

func (d *Docker) Close(serviceCtx ServiceContext) error {
	if d.client == nil {
		return nil
	}
	return d.client.Close()
}
