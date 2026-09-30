package main

import (
	"os"

	"app/bootstrap"
	"app/config"
)

func main() {
	app := bootstrap.BootApp()
	if len(os.Args) > 1 {
		os.Exit(app.HandleCommand(os.Args[1:]...))
	}
	app.Run(":" + config.App.Port)
}
