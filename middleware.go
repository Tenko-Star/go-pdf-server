package main

import (
	"github.com/kataras/iris/v12/context"
)

var (
	isStart bool
	q       chan *request
)

func processQueue(maxDepth, maxThread uint) context.Handler {
	if !isStart {
		_logger.Info("Init request queue")
		q = make(chan *request, maxDepth)
		for i := 0; i < int(maxThread); i++ {
			go PdfProcessor(q)
		}

		isStart = true
	}

	return func(c context.Context) {
		_logger.Debug("Inject request queue")
		c.Values().Set("requestQueue", q)
		c.Next()
	}
}
