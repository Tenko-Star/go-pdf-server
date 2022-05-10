package main

import (
	"fmt"
	"github.com/kataras/golog"
	"os"
	"pdf-server/config"
	"strings"
	"time"
)

var (
	globalLogger *golog.Logger
	globalConfig *config.Config

	currentWriter *logWriter
)

type logWriter struct {
	target *os.File
}

func (w *logWriter) Write(buf []byte) (int, error) {
	return w.target.Write(buf)
}

func timing() {
	var now, next time.Time
	var timer *time.Timer
	for {
		now = time.Now()
		next = now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

		timer = time.NewTimer(next.Sub(now))
		<-timer.C

		now = time.Now()
		// ->.../pdf-220428.log
		changeOutput(fmt.Sprintf("%s/pdf-%s.log", strings.Trim(globalConfig.Roots.Log, "/\\"), now.Format("060102")))
	}
}

func changeOutput(path string) {
	newer, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, os.ModePerm)

	// 创建文件出错但是旧的文件可用则继续使用旧文件
	if err != nil && currentWriter.target != nil {
		fmt.Printf("[WARN] log file could not be created. Use old one. file: %s\n\t\tBecause %s", currentWriter.target.Name(), err.Error())
		globalLogger.Warnf("log file could not be created. Use this file. Because %s", err.Error())
		return
	}
	// 创建文件出错且旧的文件不可用则直接停止运行
	if err != nil && currentWriter.target == nil {
		panic(fmt.Sprintf("[ERRO] log service could not run. Because %s", err.Error()))
	}

	temp := &logWriter{
		target: newer,
	}

	// 设置新的输出
	globalLogger.SetOutput(temp)

	if currentWriter != nil && currentWriter.target != nil {
		// 关闭旧的文件句柄
		err = currentWriter.target.Close()
		if err != nil {
			fmt.Println("[WARN] log file could not be closed.")
		}
	}

	// 保存新的输出文件信息
	currentWriter = temp
}

func RedirectLog(logger *golog.Logger) {
	var err error
	globalLogger = logger
	globalConfig = config.Get()

	if _, err = os.Stat(globalConfig.Roots.Log); err != nil {
		err = os.MkdirAll(globalConfig.Roots.Log, os.ModePerm)
	}
	if err != nil {
		logger.Errorf("Could not create log dir. Because %s", err.Error())
		panic(fmt.Sprintf("Could not create log dir. Because %s", err.Error()))
	}

	now := time.Now()
	changeOutput(fmt.Sprintf("%s/pdf-%s.log", strings.Trim(globalConfig.Roots.Log, "/\\"), now.Format("060102")))

	go timing()
}
