package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/magomed066/url-checker/internal/features/monitors"
	"github.com/magomed066/url-checker/internal/store"
)

type application struct {
	addr   string
	config config
	store  *store.Store
	client *http.Client
}

type config struct {
	dns string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.ClientIPFromRemoteAddr) // pick one ClientIPFrom* based on your infra, see below
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Server is running"))
	})

	monitorsService := monitors.NewService(app.store, app.client)
	monitorsHandler := monitors.NewHandler(monitorsService)

	r.Route("/monitors", func(r chi.Router) {
		r.Post("/", monitorsHandler.Create)
	})

	return r
}

func (app *application) run(ctx context.Context, h http.Handler) error {
	server := http.Server{
		Addr:         app.addr,
		Handler:      h,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 30,
		IdleTimeout:  time.Minute,
	}

	errorCh := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()

		errorCh <- err
	}()

	select {
	case err := <-errorCh:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			5*time.Second,
		)
		defer cancel()

		shutdownErr := server.Shutdown(shutdownCtx)
		if shutdownErr != nil {
			_ = server.Close()
		}

		serverErr := <-errorCh
		if errors.Is(serverErr, http.ErrServerClosed) {
			serverErr = nil
		}

		return errors.Join(shutdownErr, serverErr)
	}
}
