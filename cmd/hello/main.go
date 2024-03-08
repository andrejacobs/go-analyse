package main

import (
	"fmt"

	"github.com/andrejacobs/go-analyse/internal/compiledinfo"
)

func main() {
	fmt.Println("// " + compiledinfo.UsageNameAndVersion())
}
