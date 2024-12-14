package pdcs

import (
	"context"
	"log"
	"log/slog"

	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/libpod/events"
	"github.com/containers/podman/v5/pkg/domain/entities"
)

func StartEventStreamer(logger *slog.Logger, updateImage bool) {
	var eventOptions entities.EventsOptions

	logger.Info("starting podman event streamer")

	if updateImage {
		logger.Debug("update images")
		updateImages()
	}

	eventChannel := make(chan events.ReadResult, 1)
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
					logger.Error("podman received event not ok")

					continue
				}

				if updateImage && event.Event.Type == events.Image {
					logger.Debug("update images")
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
