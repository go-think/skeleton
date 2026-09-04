package main

import (
	"app/bootstrap"
	"app/config"
)

func main() {
	app := bootstrap.BootApp()
	app.Run(":" + config.App.Port)
}
