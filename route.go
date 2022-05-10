package main

import "github.com/kataras/iris/v12"

func PdfHandle(c iris.Context) {
	var (
		ok  bool
		err error

		r = &request{}
		q chan *request
	)

	if err = c.ReadJSON(r); err != nil {
		c.Application().Logger().Infof("Incorrect request body format.\n[Data]\n%s\n[End]\nIf data is empty, maybe request body is incorrect.", getRequestContent(c))
		c.StatusCode(400)
		_, _ = c.JSON(iris.Map{
			"code": 0,
			"msg":  "Bad Request",
			"data": nil,
		})
		return
	}

	if q, ok = c.Values().Get("requestQueue").(chan *request); !ok {
		c.Application().Logger().Errorf("Could not get request queue from context, please check code.")
		c.StatusCode(500)
		_, _ = c.JSON(iris.Map{
			"code": 0,
			"msg":  "Internal error",
			"data": nil,
		})

		c.StopExecution()
		return
	}

	if cap(q) == len(q) {
		// 队列已满
		c.Application().Logger().Warnf("Could not to process more request, because queue is full.")
		c.StatusCode(503)
		_, _ = c.JSON(iris.Map{
			"code": 0,
			"msg":  "Service Unavailable",
			"data": nil,
		})

		c.StopExecution()
		return
	}

	q <- r
	c.Application().Logger().Debugf("Add new request.")

	_, _ = c.JSON(iris.Map{
		"code": 0,
		"msg":  "success",
		"data": nil,
	})
}

func ProcessorHandle(c iris.Context) {
	var (
		err error
	)

	_, err = c.JSON(iris.Map{
		"code": 0,
		"msg":  "ok",
		"data": map[string]interface{}{
			"in_process": getInProcessNumber(),
			"success":    getSuccessNumber(),
			"failure":    getFailureNumber(),
		},
	})
	if err != nil {
		_logger.Warnf("Response error: %s", err.Error())
	}
}
