package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moukoublen/goboilerplate/build"
)

func main() {
	fmt.Printf(
		"Starting (version %s)(branch %s)(commit %s)(commit short %s)(tag %s)\n",
		build.Version,
		build.Branch,
		build.Commit,
		build.CommitShort,
		build.Tag,
	)

	ctx, cancel := context.WithCancel(context.Background())
	_ = ctx

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalCh
	fmt.Printf("Signal received: %d %s\n", sig, sig.String())
	cancel()
}
