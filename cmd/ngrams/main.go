package main

import (
	"errors"
	"os"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
)

func main() {
	if err := app.Main(os.Stdout, os.Stderr); err != nil {
		die(err, 1)
	}
}

func die(err error, code int) {
	if errors.Is(err, app.ErrExitWithNoErr) {
		os.Exit(0)
	}

	os.Exit(code)
}
