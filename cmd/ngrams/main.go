package main

import (
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

	_ = appCmd
}

func die(err error, code int) {
	fmt.Fprintf(os.Stderr, "%v\n", err)
	os.Exit(code)
}
