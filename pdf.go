package main

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/kataras/iris/v12"
	"hash"
	"io"
	"os"
	"pdf-server/http"
	"pdf-server/pdf"
	"sync"
	"time"
)

type pdfGuard struct {
	id        uint64
	inProcess uint64
	success   uint64
	failure   uint64
}

var guard pdfGuard
var lock sync.Mutex

func PdfProcessor(q chan *request) {
	var (
		cid uint64
		err error

		r        *request
		fName    string
		p        = pdf.New(_logger)
		response []byte
	)

	for {
		r = <-q

		cid = getCid()

		_logger.Infof("CID: %v PID: %v\n\t\t\t\t\t\tStart process pdf file.", cid, r.Id)

		if r.FileName == "" {
			fName = getPdfFileName(r.Id)
		} else {
			fName = getPdfFileName(r.FileName)
		}

		p.SetOutput(fName)

		if r.Raw != nil {
			p.FromData(r.Raw)
		} else if r.Template != nil {
			err = p.FromTemplate(r.Template.Path, r.Template.Replace, r.Template.Insert)
			if err != nil {
				_logger.Debugf("FromTemplate error")
				setFailure(cid, err)
				callFailure(r.Id, err.Error())
				continue
			}
		} else {
			_logger.Debugf("No data input error")
			setFailure(cid, "No data input")
			callFailure(r.Id, "No data input")
			continue
		}

		// 处理pdf
		err = p.Save()
		if err != nil {
			_logger.Debugf("Save error")
			setFailure(cid, err)
			callFailure(r.Id, err.Error())
			continue
		}

		if _, err = os.Stat(fName); err == nil {
			response, err = http.UploadFile(fName, _config.Upload, "file")
			if err != nil {
				setFailure(cid, err)
				callFailure(r.Id, err.Error())
				continue
			}
			data := struct {
				Data struct {
					Id int `json:"id"`
				}
			}{}
			err = json.Unmarshal(response, &data)
			if err != nil {
				setFailure(cid, err)
				callFailure(r.Id, err.Error())
				continue
			}

			setSuccessful(r.Id, cid)
			callSuccess(r.Id, data.Data.Id)

			if !_config.KeepTemp {
				err = os.Remove(fName)

				if err != nil {
					_logger.Warnf("Could not to delete temp pdf file, please delete temp later.")
				}
			}
		} else {
			_logger.Debugf("os.Stat error")
			setFailure(cid, err.Error())
			callFailure(r.Id, err.Error())
		}

	}
}

func getRequestContent(c iris.Context) string {
	data, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return ""
	}

	return string(data)
}

func getPdfFileName(key interface{}) string {
	var (
		err error

		oRoot  = _config.Roots.Output
		random = _config.Debug

		m    hash.Hash
		base []byte
		buf  bytes.Buffer
		enc  *gob.Encoder
	)

	if str, ok := key.(string); ok {
		return fmt.Sprintf("%s/%s.pdf", oRoot, str)
	}

	enc = gob.NewEncoder(&buf)
	err = enc.Encode(key)
	if err != nil {
		base = []byte("")
	} else {
		base = buf.Bytes()
	}

	m = md5.New()
	m.Write(base)
	if random {
		m.Write([]byte(time.Now().String()))
	}

	return fmt.Sprintf("%s/%s.pdf", oRoot, hex.EncodeToString(m.Sum(nil)))
}

type atomAction func()

func atom(action atomAction) {
	lock.Lock()
	action()
	lock.Unlock()
}

func getCid() uint64 {
	var (
		id uint64
	)

	atom(func() {
		id = guard.id
		if guard.id == 0xffffffffffffffff {
			_logger.Warnf("Number of id is out of range.")
			guard.id = 0
		} else {
			guard.id++
		}

		if guard.inProcess == 0xffffffffffffffff {
			_logger.Warnf("Number of processing events is out of range.")
			guard.inProcess = 0
		} else {
			guard.inProcess++
		}
	})

	return id
}

func setSuccessful(id int, cid uint64) {
	atom(func() {
		if guard.success == 0xffffffffffffffff {
			_logger.Warnf("Number of success is out of range.")
			guard.success = 0
		} else {
			guard.success++
		}
	})

	_logger.Infof("Succeeded to create pdf. Id: %d. Cid: %d", id, cid)
}

func setFailure(cid uint64, reason interface{}) {
	atom(func() {
		if guard.failure == 0xffffffffffffffff {
			_logger.Warnf("Number of failure is out of range.")
			guard.failure = 0
		} else {
			guard.failure++
		}
	})

	switch reason.(type) {
	case string:
		_logger.Infof("[CID:%v]Failed to create, because %s", cid, reason.(string))

	case error:
		_logger.Infof("[CID:%v]Failed to create, because %s", cid, reason.(error).Error())

	default:
		_logger.Infof("[CID:%v]Failed to create", cid)
	}
}

func getInProcessNumber() uint64 {
	return guard.inProcess
}

func getSuccessNumber() uint64 {
	return guard.success
}

func getFailureNumber() uint64 {
	return guard.failure
}
