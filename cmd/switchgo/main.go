package main

import (
	"fmt"
	"os"

	"github.com/meiyoutoufa/switchgo/pkg/cmd"
)

func main() {
	if err := cmd.GetRootCommand().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
