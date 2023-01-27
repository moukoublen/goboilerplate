package internal

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/rs/zerolog/log"
)

// ChannelForSignals returns a channel in which each one of the `signals` will be sent (if received by the process).
func ChannelForSignals(bufferSize int, signals []os.Signal) <-chan os.Signal {
	signalCh := make(chan os.Signal, bufferSize)
	signal.Notify(signalCh, signals...)

	return signalCh
}

func WaitForFatal(signalsCh <-chan os.Signal, fatalErrorsCh <-chan error) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		for {
			select {
			case sig := <-signalsCh:
				log.Info().Msgf("Signal received: %d %s", sig, sig.String())
			case servErr := <-fatalErrorsCh:
				if !errors.Is(servErr, http.ErrServerClosed) {
					log.Error().Err(servErr).Msg("fatal error received")
				}
			}
		}
	}()
}
