package internal

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Daemon struct wraps tha basic functionality that is needed in order for an application to run daemon/service like and shutdown gracefully when stop conditions are met.
// Stop conditions:
//
//	a. a signal (one of daemonConfig.signalsNotify) is received from OS.
//	b. a error is received in fatal errors channel.
//	c. the given root context (`rootCTX`) in `NewDaemon` function is done.
//
// As described in `b` a fatal error channel is being provided (function `FatalErrorsChannel()`) and can be used by the rest of the code when a catastrophic error occurs that needs to trigger an application shutdown.
type Daemon struct {
	signalCh        chan os.Signal
	fatalErrorsCh   chan error
	done            chan struct{}
	logger          *slog.Logger
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

func (o *Daemon) shutDown(ctx context.Context) {
	o.logger.InfoContext(ctx, "starting graceful shutdown")
	deadline := time.Now().Add(o.config.shutdownTimeout)
	dlCtx, dlCancel := context.WithDeadline(ctx, deadline)
	o.onShutDownMutex.Lock()
	for _, f := range o.onShutDown {
		f(dlCtx)
	}
	o.onShutDownMutex.Unlock()
	dlCancel()

	close(o.fatalErrorsCh)
	signal.Stop(o.signalCh)
	close(o.done)
}

// FatalErrorsChannel returns the fatal error channel that can be used by the application in order to trigger a shutdown.
func (o *Daemon) FatalErrorsChannel() chan<- error {
	return o.fatalErrorsCh
}

// start will spawn a go routine that will run until one of the stop conditions is met.
// After a stop conditions is met the `Daemon` will attempt shutdown "gracefully" by running every function that is registered in `onShutDown` slice, sequentially.
func (o *Daemon) start(rootCTX context.Context) context.Context {
	ctx, cnl := context.WithCancel(rootCTX)

	// consume fatal errors
	go func() {
		for err := range o.fatalErrorsCh {
			// Stop condition (B) fatal error received.
			logFatalErr(o.logger, err)
			cnl()
		}
	}()

	// consume signals
	go func() {
		sigReceived := 0
		for sig := range o.signalCh {
			// Stop condition (A) signal received.
			sigReceived++
			logSig(o.logger, sig)
			cnl()
			if sigReceived == o.config.maxSignalCount {
				o.logger.Error("max number of signal received. Terminating immediately")
				os.Exit(1) //nolint:revive //its the emergency shut-down situation
			}
		}
	}()

	go func() {
		select {
		// Stop condition (C) root context is done.
		case <-rootCTX.Done():
			o.logger.Error("root context got canceled", slog.String("error", rootCTX.Err().Error()))
			cnl()
			o.shutDown(context.Background()) //nolint:contextcheck

			return

		case <-ctx.Done():
			if err := rootCTX.Err(); err != nil {
				o.logger.Error("context got canceled", slog.String("error", rootCTX.Err().Error()))
			} else {
				o.logger.Error("context got canceled")
			}

			o.shutDown(rootCTX)
		}
	}()

	return ctx
}

func (o *Daemon) Wait() {
	<-o.done
	o.logger.Info("shutdown completed")
}

type DaemonConfigOption func(*daemonConfig)

// SetSignalsNotify sets the OS signals that will be used as stop condition to Daemon in order to shutdown gracefully.
func SetSignalsNotify(signals ...os.Signal) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.signalsNotify = signals
	}
}

// SetMaxSignalCount sets the maximum number of signals to receive while waiting for graceful shutdown.
// If the max number of signals exceeds, immediate termination will follow.
func SetMaxSignalCount(size int) DaemonConfigOption {
	return func(oc *daemonConfig) {
		oc.maxSignalCount = size
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
	defaultMaxSignalCount               = 2
	defaultFatalErrorsChannelBufferSize = 10
	defaultShutdownTimeout              = 4 * time.Second
)

type daemonConfig struct {
	signalsNotify                []os.Signal
	maxSignalCount               int
	fatalErrorsChannelBufferSize int
	shutdownTimeout              time.Duration
}

func NewDaemon(ctx context.Context, logger *slog.Logger, opts ...DaemonConfigOption) (*Daemon, context.Context) {
	cnf := daemonConfig{
		signalsNotify:                []os.Signal{os.Interrupt, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM},
		maxSignalCount:               defaultMaxSignalCount,
		fatalErrorsChannelBufferSize: defaultFatalErrorsChannelBufferSize,
		shutdownTimeout:              defaultShutdownTimeout,
	}

	for _, o := range opts {
		o(&cnf)
	}

	signalCh := make(chan os.Signal, cnf.maxSignalCount)
	signal.Notify(signalCh, cnf.signalsNotify...)

	o := &Daemon{
		signalCh:      signalCh,
		fatalErrorsCh: make(chan error, cnf.fatalErrorsChannelBufferSize),
		done:          make(chan struct{}),
		config:        cnf,
		logger:        logger,
	}

	return o, o.start(ctx)
}

func logSig(logger *slog.Logger, sig os.Signal) {
	var signalStr string
	if sigInt, ok := sig.(syscall.Signal); ok {
		signalStr = fmt.Sprintf("%s %d", sig.String(), int(sigInt))
	} else {
		signalStr = sig.String()
	}

	logger.Warn("signal received", slog.String("signal", signalStr))
}

func logFatalErr(logger *slog.Logger, err error) {
	logger.Error("fatal error received", slog.String("error", err.Error()))
}
