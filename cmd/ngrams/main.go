package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
)

func main() {
	opts, err := app.ParseArgs()
	if err != nil {
		die(err, 1)
	}

	appCmd, err := app.New(opts...)
	if err != nil {
		die(err, 1)
	}

	if err := appCmd.Run(os.Stdout); err != nil {
		die(err, 2)
	}
}

func die(err error, code int) {
	if errors.Is(err, app.ErrExitWithNoErr) {
		os.Exit(0)
	}

	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(code)
}
