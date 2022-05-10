package main

import "pdf-server/http"

func callSuccess(id int, file int) {
	var (
		err      error
		response []byte

		callback = _config.Callback.Success
	)

	response, err = http.PostWithJson(callback, map[string]interface{}{
		"id":   id,
		"file": file,
	}, nil)

	if err == nil {
		_logger.Infof("Success To Callback: %v.", id)
	} else {
		_logger.Infof("Fail To Callback: %v.", id)
	}

	if _config.Debug && err != nil {
		_logger.Debugf("[Success callback] PID: %d\n\t\t\t\tRes: \t%s\n\t\t\t\tErr: \t%s", id, string(response), err.Error())
	} else {
		_logger.Debugf("[Success callback] PID: %d\n\t\t\t\tRes: \t%s", id, string(response))
	}
}

func callFailure(id int, reason string) {
	var (
		err      error
		response []byte

		callback = _config.Callback.Success
	)

	response, err = http.PostWithJson(callback, map[string]interface{}{
		"id":     id,
		"reason": reason,
	}, nil)

	if err == nil {
		_logger.Infof("Success To Callback: %v.", id)
	} else {
		_logger.Infof("Fail To Callback: %v.", id)
	}

	if _config.Debug && err != nil {
		_logger.Debugf("[Fail callback] PID: %d\n\t\t\t\tRes: \t%s\n\t\t\t\tErr: \t%s", id, string(response), err.Error())
	} else {
		_logger.Debugf("[Fail callback] PID: %d\n\t\t\t\tRes: \t%s", id, string(response))
	}
}
