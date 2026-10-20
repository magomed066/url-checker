package main

import "os"

func main() {
	app := application{
		addr: ":8080",
	}

	if err := app.run(app.mount()); err != nil {
		os.Exit(1)
	}
}
