package app

import (
	"context"
	"net/http"
	"time"
)

type App struct {
	Context context.Context
}

// New returns a new instance of the App
func New(ctx context.Context) (*App, error) {
	return &App{
		Context: ctx,
	}, nil
}

func (a *App) Serve(address string) error {
	srv := http.Server{
		Addr:         address,
		Handler:      nil,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		<-a.Context.Done()
		srv.Close()
	}()

	err := srv.ListenAndServe()
	if err == http.ErrServerClosed {
		return nil
	}

	return err
}
