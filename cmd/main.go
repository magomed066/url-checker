package main

import (
	"os"

	"github.com/magomed066/url-checker/internal/store"
)

func main() {

	store := store.NewStore()

	app := application{
		addr:  ":8080",
		store: store,
	}

	if err := app.run(app.mount()); err != nil {
		os.Exit(1)
	}
}
