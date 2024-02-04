package internal

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

// Daemon struct wraps tha basic functionality that is needed in order for an application to run daemon/service like and shutdown gracefully when stop conditions are met.
// Stop conditions:
//
//	a. a signal (one of daemonConfig.signalsNotify) is received from OS.
//	b. a error is received in fatal errors channel.
//	c. the given context (`ctx`) in `.Daemon` function is done.
//
// As described in `b` a fatal error channel is being provided (function `FatalErrorsChannel()`) and can be used by the rest of the code when a catastrophic error occurs that needs to trigger an application shutdown.
type Daemon struct {
	signalCh        chan os.Signal
	fatalErrorsCh   chan error
	onShutDown      []func(context.Context)
	config          daemonConfig
	onShutDownMutex sync.Mutex
}

// OnShutDown appends a function to be called on shutdown.
func (o *Daemon) OnShutDown(f ...func(context.Context)) {
	o.onShutDownMutex.Lock()
	defer o.onShutDownMutex.Unlock()
	o.onShutDown = append(o.onShutDown, f...)
}

func (o *Daemon) callOnShutDown(ctx context.Context) {
	o.onShutDownMutex.Lock()
	defer o.onShutDownMutex.Unlock()
	for _, f := range o.onShutDown {
		f(ctx)
	}
}

// FatalErrorsChannel returns the fatal error channel that can be used by the application in order to trigger a shutdown.
func (o *Daemon) FatalErrorsChannel() chan<- error {
	return o.fatalErrorsCh
}

// Run will block until one of the stop conditions is met.
// After a stop conditions is met the `Daemon` will attempt shutdown "gracefully" by running every function that is registered in `onShutDown` slice, sequentially.
func (o *Daemon) Run(ctx context.Context, ctxCancel context.CancelFunc) {
	log := zerolog.Ctx(ctx)

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

type DaemonConfigOption func(*daemonConfig)

// SetSignalsNotify sets the OS signals that will be used as stop condition to Daemon in order to shutdown gracefully.
func SetSignalsNotify(signals ...os.Signal) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.signalsNotify = signals
	}
}

// SetSignalsChannelBufferSize sets the channel size of the watched received signals in case that is needed to be a buffered one.
func SetSignalsChannelBufferSize(size int) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.signalChannelBufferSize = size
	}
}

// SetFatalErrorsChannelBufferSize sets the fatal error channel size in case that is needed to be a buffered one.
func SetFatalErrorsChannelBufferSize(size int) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.fatalErrorsChannelBufferSize = size
	}
}

// SetShutdownTimeout sets a timeout to the graceful shutdown process.
func SetShutdownTimeout(d time.Duration) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.shutdownTimeout = d
	}
}

const (
	defaultSignalChannelBufferSize      = 10
	defaultFatalErrorsChannelBufferSize = 10
	defaultShutdownTimeout              = 4 * time.Second
)

type daemonConfig struct {
	signalsNotify                []os.Signal
	signalChannelBufferSize      int
	fatalErrorsChannelBufferSize int
	shutdownTimeout              time.Duration
}

func NewDaemon(opts ...DaemonConfigOption) *Daemon {
	cnf := daemonConfig{
		signalsNotify:                []os.Signal{os.Interrupt, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM},
		signalChannelBufferSize:      defaultSignalChannelBufferSize,
		fatalErrorsChannelBufferSize: defaultFatalErrorsChannelBufferSize,
		shutdownTimeout:              defaultShutdownTimeout,
	}

	for _, o := range opts {
		o(&cnf)
	}

	signalCh := make(chan os.Signal, cnf.signalChannelBufferSize)
	signal.Notify(signalCh, cnf.signalsNotify...)

	o := &Daemon{
		signalCh:      signalCh,
		fatalErrorsCh: make(chan error, cnf.fatalErrorsChannelBufferSize),
		config:        cnf,
	}

	return o
}
