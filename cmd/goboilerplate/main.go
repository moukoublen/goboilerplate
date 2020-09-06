package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/moukoublen/goboilerplate/internal"
)

func main() {
	fmt.Printf(
		"Starting %s service (version %s)(branch %s)(commit %s)(commit short %s)(tag %s)\n",
		internal.Name,
		internal.Version,
		internal.Branch,
		internal.Commit,
		internal.CommitShort,
		internal.Tag,
	)

	ctx, cnl := context.WithCancel(context.Background())
	_ = ctx

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	sig := <-signalCh
	fmt.Printf("Signal received: %d %s\n", sig, sig.String())
	cnl()
}
