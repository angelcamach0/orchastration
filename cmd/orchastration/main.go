package main

import (
	"fmt"
	"os"

	"orchastration/internal/app"
	"orchastration/internal/version"
)

func main() {
	code, err := app.Run(os.Args[1:], version.Info())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(code)
}
