package main

import (
	"net/http"
	"os"
	"time"

	"github.com/magomed066/url-checker/internal/store"
)

func main() {

	store := store.NewStore()

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	app := application{
		addr:   ":8080",
		store:  store,
		client: client,
	}

	if err := app.run(app.mount()); err != nil {
		os.Exit(1)
	}
}
