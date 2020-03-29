package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/i25959341/sku-aggregator/internal/app"
	"github.com/i25959341/sku-aggregator/internal/config"
)

func New() error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go shutdownHandler(cancelFunc)

	conf := config.New()

	a, err := app.New(ctx)
	if err != nil {
		return err
	}

	return a.Serve(fmt.Sprintf(":%v", conf.HTTPPort))
}

// shutdownHandler listens for a SIGTERM signal
// and gracefully cancels the main application context
// once this is completed exits the app
func shutdownHandler(cancelFunction context.CancelFunc) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-sig

	// Invoke the cancel function
	cancelFunction()

	// Safety deadline for exiting
	<-time.After(10 * time.Second)
	os.Exit(1)
}
