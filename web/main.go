package main

import (
	"os"

	"./router"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/AmyangXYZ/sweetygo/middlewares"
)

var (
	addr   = ":8080"
	tplDir = "templates"
)

func main() {
	app := sweetygo.New()
	app.SetTemplates(tplDir, nil)
	app.USE(middlewares.Logger(os.Stdout, middlewares.DefaultSkipper))

	router.SetRouter(app)

	app.Run(":8080")
}
