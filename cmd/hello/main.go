package main

import (
	"fmt"

	"github.com/andrejacobs/go-analysis/internal/compiledinfo"
)

func main() {
	fmt.Println("// " + compiledinfo.UsageNameAndVersion())
}
