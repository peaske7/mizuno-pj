package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

const (
	version = "0.1.0"
	appname = "adelaida-cli"
)

func main() {
	err := newApp().Run(os.Args)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func newApp() *cli.App {
	app := &cli.App{
		Name:                 appname,
		Version:              version,
		EnableBashCompletion: true,
		Commands:             commands,
	}

	return app
}

var commands = []*cli.Command{
	convert,
}
