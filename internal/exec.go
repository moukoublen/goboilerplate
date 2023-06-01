package internal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

type MainConfigOption func(*mainConfig)

func SetSignalsNotify(signals ...os.Signal) MainConfigOption {
	return func(oc *mainConfig) {
		oc.signalsNotify = signals
	}
}

func SetSignalsChannelBufferSize(size int) MainConfigOption {
	return func(oc *mainConfig) {
		oc.signalChannelBufferSize = size
	}
}

func SetFatalErrorsChannelBufferSize(size int) MainConfigOption {
	return func(oc *mainConfig) {
		oc.fatalErrorsChannelBufferSize = size
	}
}

func SetShutdownTimeout(d time.Duration) MainConfigOption {
	return func(oc *mainConfig) {
		oc.shutdownTimeout = d
	}
}

const (
	defaultSignalChannelBufferSize      = 10
	defaultFatalErrorsChannelBufferSize = 10
	defaultShutdownTimeout              = 4 * time.Second
)

type mainConfig struct {
	signalsNotify                []os.Signal
	signalChannelBufferSize      int
	fatalErrorsChannelBufferSize int
	shutdownTimeout              time.Duration
}

func NewMain(opts ...MainConfigOption) *Main {
	cnf := mainConfig{
		signalsNotify:                []os.Signal{os.Interrupt, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM},
		signalChannelBufferSize:      defaultSignalChannelBufferSize,
		fatalErrorsChannelBufferSize: defaultFatalErrorsChannelBufferSize,
		shutdownTimeout:              defaultShutdownTimeout,
	}

	for _, o := range opts {
		o(&cnf)
	}

	o := &Main{
		signalCh:      createSignalsChannel(cnf),
		fatalErrorsCh: make(chan error, cnf.fatalErrorsChannelBufferSize),
		config:        cnf,
	}

	return o
}

type Main struct {
	onShutDownMutex sync.Mutex
	onShutDown      []func(context.Context)

	signalCh      chan os.Signal
	fatalErrorsCh chan error
	config        mainConfig
}

// OnShutDown appends a function to be called on shutdown.
func (o *Main) OnShutDown(f func(context.Context)) {
	o.onShutDownMutex.Lock()
	defer o.onShutDownMutex.Unlock()
	o.onShutDown = append(o.onShutDown, f)
}

func (o *Main) callOnShutDown(ctx context.Context) {
	o.onShutDownMutex.Lock()
	defer o.onShutDownMutex.Unlock()
	for _, f := range o.onShutDown {
		f(ctx)
	}
}

func (o *Main) FatalErrorsChannel() chan<- error {
	return o.fatalErrorsCh
}

func (o *Main) Run(ctx context.Context, ctxCancel context.CancelFunc) {
	logSig := func(sig os.Signal) {
		event := log.Warn().Str("signal", sig.String())
		if sigInt, ok := sig.(syscall.Signal); ok {
			event.Int("signalNumber", int(sigInt))
		}
		event.Msg("signal received")
	}
	logFatalErr := func(err error) {
		log.Error().Err(err).Msg("fatal error received")
	}

	select {
	case sig := <-o.signalCh:
		logSig(sig)
		go func() {
			for sig := range o.signalCh {
				logSig(sig)
			}
		}()
	case fatalErr := <-o.fatalErrorsCh:
		logFatalErr(fatalErr)
		go func() {
			for err := range o.fatalErrorsCh {
				logFatalErr(err)
			}
		}()
	case <-ctx.Done():
		log.Error().Err(ctx.Err()).Msg("root context got canceled")
	}

	// graceful shut down.
	log.Info().Msgf("starting graceful shutdown")
	deadline := time.Now().Add(o.config.shutdownTimeout)
	dlCtx, dlCancel := context.WithDeadline(ctx, deadline)

	o.callOnShutDown(dlCtx)

	dlCancel()
	ctxCancel()
	close(o.fatalErrorsCh)
	log.Info().Msgf("shutdown completed")
}

func createSignalsChannel(cnf mainConfig) chan os.Signal {
	signalCh := make(chan os.Signal, cnf.signalChannelBufferSize)
	signal.Notify(signalCh, cnf.signalsNotify...)

	return signalCh
}
