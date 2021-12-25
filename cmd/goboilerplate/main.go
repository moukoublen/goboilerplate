package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moukoublen/goboilerplate/internal"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	router := internal.SetupDefaultRouter()
	server, chErr := internal.StartListenAndServe(":43000", router)
	go func() {
		err := <-chErr
		if !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("http server error: %s\n", err.Error())
		}
	}()

	blockForSignals(os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	// shutdown gracefully the http server.
	dlCtx, dlCancel := context.WithDeadline(ctx, time.Now().Add(3*time.Second))
	_ = server.Shutdown(dlCtx)
	dlCancel()

	cancel()
}

func blockForSignals(s ...os.Signal) {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, s...)
	sig := <-signalCh
	fmt.Printf("Signal received: %d %s\n", sig, sig.String())
	close(signalCh)
}
