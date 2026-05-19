package main

import (
	"fmt"
	"os"

	"github.com/cryptdefender3232/phantom/util/assets"
)

func main() {
	if err := assets.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
