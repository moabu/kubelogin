package main

import (
	"os"

	"github.com/moabu/kubelogin/pkg/cmd"
)

func main() {

	root := cmd.NewRootCmd()
	if err := root.Execute(); err != nil {
		os.Exit(1)
	}

}
