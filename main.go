package main

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"pdf-server/config"
)

var (
	_config = config.Get()
	_logger *golog.Logger
)

func main() {
	app := iris.Default()
	_logger = app.Logger()

	appInit()

	appUseMiddles(app)

	appUseRoutes(app)

	err := app.Run(iris.Addr(fmt.Sprintf("%s:%s", _config.Server.Host, _config.Server.Port)))
	if err != nil {
		app.Logger().Errorf("%s\n", err.Error())
		return
	}
}

func appInit() {
	RedirectLog(_logger)

	if _config.Debug {
		_logger.SetLevel("debug")
	} else {
		_logger.SetLevel("info")
	}
}

func appUseMiddles(app *iris.Application) {
	app.Use(processQueue(_config.MaxQueue, _config.MaxThread))
}

func appUseRoutes(app *iris.Application) {
	app.Post("/pdf", PdfHandle)
	app.Get("/info", ProcessorHandle)
}
