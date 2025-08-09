package pdcs

import (
	"context"
	"log"
	"log/slog"

	"github.com/containers/podman/v5/cmd/podman/registry"
	"github.com/containers/podman/v5/libpod/events"
	"github.com/containers/podman/v5/pkg/domain/entities"
)

func StartEventStreamer(logger *slog.Logger, updateImage bool) { //nolint:cyclop
	var eventOptions entities.EventsOptions

	logger.Info("starting podman event streamer")

	if updateImage {
		logger.Debug("update images")
		updateImages()
	}

	errChannel := make(chan error)
	restartChannel := make(chan bool)

	for {
		eventChannel := make(chan events.ReadResult, 1)
		eventOptions.EventChan = eventChannel
		eventOptions.Stream = true
		eventOptions.Filter = []string{}

		go func() {
			logger.Debug("starting engine event reader")

			err := registry.ContainerEngine().Events(context.Background(), eventOptions)
			if err != nil {
				errChannel <- err
			}
		}()

		go func() {
			logger.Debug("starting event reader loop")

			for {
				select {
				case event, ok := <-eventChannel:
					if !ok {
						logger.Error("podman received event not ok")

						restartChannel <- true

						return
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

		<-restartChannel
	}
}
