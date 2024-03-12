package main

import (
	"fmt"
	"os"

	"github.com/andrejacobs/go-analyse/cmd/ngrams/app"
	"github.com/andrejacobs/go-analyse/internal/compiledinfo"
)

func main() {
	fmt.Println("// " + compiledinfo.UsageNameAndVersion())

	appCmd, err := app.New(app.WithDefaults())
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	_ = appCmd
}
