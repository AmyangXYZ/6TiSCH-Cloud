package main

import (
	"os"

	"./router"
	"github.com/AmyangXYZ/sweetygo"
	"github.com/AmyangXYZ/sweetygo/middlewares"
)

var (
	addr   = ":8888"
	tplDir = "templates"
)

func main() {
	app := sweetygo.New()

	app.SetTemplates(tplDir, nil)
	app.USE(middlewares.Logger(os.Stdout, middlewares.DefaultSkipper))
	app.USE(middlewares.CORS(middlewares.CORSOpt{}))
	router.SetRouter(app)

	app.Run(addr)
}
