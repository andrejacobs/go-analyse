package main

import (
	"fmt"
	"os"

	"github.com/andrejacobs/go-analysis/internal/compiledinfo"
	"github.com/andrejacobs/go-analysis/internal/helloworld"
)

func main() {
	fmt.Println("// " + compiledinfo.UsageNameAndVersion())

	if err := helloworld.SayHello(os.Stdout, "world"); err != nil {
		fmt.Fprintf(os.Stderr, "failed to say hello. %v", err)
		os.Exit(1)
	}
}
