package pdcs

import (
	"context"
	"log"

	"github.com/containers/podman/v4/cmd/podman/registry"
	"github.com/containers/podman/v4/libpod/events"
	"github.com/containers/podman/v4/pkg/domain/entities"
	klog "github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func StartEventStreamer(logger klog.Logger) {
	var eventOptions entities.EventsOptions

	level.Debug(logger).Log("msg", "starting podman event streamer")
	level.Debug(logger).Log("msg", "update images")
	updateImages()

	eventChannel := make(chan *events.Event, 1)
	eventOptions.EventChan = eventChannel
	eventOptions.Stream = true
	eventOptions.Filter = []string{}
	errChannel := make(chan error)

	go func() {
		err := registry.ContainerEngine().Events(context.Background(), eventOptions)
		if err != nil {
			errChannel <- err
		}
	}()

	go func() {
		for {
			select {
			case event, ok := <-eventChannel:
				if !ok {
					level.Error(logger).Log("msg", "podman received event not ok")

					continue
				}

				if event.Type == events.Image {
					level.Debug(logger).Log("msg", "update images")
					updateImages()
				}
			case err := <-errChannel:
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
}
