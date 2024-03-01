package helloworld

import (
	"fmt"
	"io"
)

// SayHello writes a greeting to the provided io.Writer.
func SayHello(w io.Writer, name string) error {
	_, err := fmt.Fprintf(w, "Hello %s, it is so nice to meet you.\n", name)
	return err
}
