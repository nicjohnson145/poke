package main

import (
	"fmt"
	"os"

	"github.com/nicjohnson145/poke/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
