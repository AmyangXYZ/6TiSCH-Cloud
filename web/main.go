package main

import (
	"time"

	"github.com/AmyangXYZ/sgo"
	"github.com/AmyangXYZ/sgo/middlewares"
)

var (
	addr               = ":80"
	tplDir             = "templates"
	lastBootTime int64 = 0
)

func main() {
	app := sgo.New()

	go func() {
		for {
			lastBootTime = modelGetLastBootTime()
			time.Sleep(1 * time.Minute)
		}
	}()

	app.SetTemplates(tplDir, nil)
	// app.USE(middlewares.Logger(os.Stdout, middlewares.DefaultSkipper))
	app.USE(middlewares.CORS(middlewares.CORSOpt{}))
	SetRouter(app)

	app.Run(addr)
}
