package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/magomed066/url-checker/internal/store"
)

func main() {
	// ShutDown
	ctx := context.Background()

	// trap Ctrl+C and call cancel on the context
	ctx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	defer func() {
		signal.Stop(c)
		cancel()
	}()

	go func() {
		select {
		case <-c:
			cancel()

		case <-ctx.Done():
		}
	}()

	store := store.NewStore()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	app := application{
		addr:   ":8080",
		store:  store,
		client: client,
	}

	if err := app.run(ctx, app.mount()); err != nil {
		os.Exit(1)
	}
}
