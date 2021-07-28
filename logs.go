package go_logs

import (
	"fmt"
	"github.com/kpango/glg"
	"io"
	"log"
	"os"
	"runtime"
	"strings"
	"syscall"
	"time"
)

const (
	DEBUG = 1
	INFO  = 2
	WARN  = 4
	ERR   = 8
)

// GinLog gin.DefaultWriter = io.MultiWriter(inits.GinLog, os.Stdout)
var GinLog = glg.FileWriter("logs/gin.log", 0777)
var infoLog = glg.FileWriter("logs/info.log", 0777)
var errLog = glg.FileWriter("logs/error.log", 0777)
var debugLog = glg.FileWriter("logs/debug.log", 0777)
var warnLog = glg.FileWriter("logs/warn.log", 0777)
var WRITE = byte(INFO | ERR)
var StderrFile = false
var stdErrFileHandler *os.File
var Err = glg.Error

func init() {
	glg.Get().
		SetMode(glg.BOTH) // default is STD
	// DisableColor().
	// SetMode(glg.NONE).
	// SetMode(glg.WRITER).
	// SetMode(glg.BOTH).
	// InitWriter().
	// AddWriter(customWriter).
	// SetWriter(customWriter).
	// AddLevelWriter(glg.LOG, customWriter).
	// AddLevelWriter(glg.INFO, customWriter).
	// AddLevelWriter(glg.WARN, customWriter).
	// AddLevelWriter(glg.ERR, customWriter).
	// SetLevelWriter(glg.LOG, customWriter).
	// SetLevelWriter(glg.INFO, customWriter).
	// SetLevelWriter(glg.WARN, customWriter).
	// SetLevelWriter(glg.ERR, customWriter).
	go splitLogByDay()
	go rewriteStderrFile()
}
func splitLogByDay() {
	timeDay := "2006-01-02"
	for {
		NowTimeDay := time.Now().Format("2006-01-02")
		if NowTimeDay > timeDay {
			if (WRITE & INFO) != 0 {
				glg.Get().SetLevelWriter(glg.INFO, io.MultiWriter(glg.FileWriter(fmt.Sprintf("logs/%s-info.log", NowTimeDay), 0777), infoLog))
			}
			if (WRITE & ERR) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("logs/%s-err.log", NowTimeDay), 0777), errLog))
			}
			if (WRITE & WARN) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("logs/%s-warn.log", NowTimeDay), 0777), warnLog))
			}
			if (WRITE & DEBUG) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("logs/%s-debug.log", NowTimeDay), 0777), debugLog))
			}
		}
		time.Sleep(time.Second)
	}
}

func rewriteStderrFile() {
	if runtime.GOOS == "windows" || !StderrFile {
		return
	}
	file, err := os.OpenFile("logs/stdErr.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		Err(err)
		return
	}
	stdErrFileHandler = file //把文件句柄保存到全局变量，避免被GC回收
	if err = syscall.Dup2(int(file.Fd()), int(os.Stderr.Fd())); err != nil {
		Err(err)
		return
	}
	// 内存回收前关闭文件描述符
	runtime.SetFinalizer(stdErrFileHandler, func(fd *os.File) {
		fd.Close()
	})
	return
}

func Debug(val ...interface{}) {
	val = append([]interface{}{findCaller(0)}, val...)
	err := glg.Debug(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Info(val ...interface{}) {
	val = append([]interface{}{findCaller(0)}, val...)
	err := glg.Info(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Warn(val ...interface{}) {
	val = append([]interface{}{findCaller(0)}, val...)
	err := glg.Warn(val...)
	if err != nil {
		log.Println(err)
		return
	}
}

//func Err(val ...interface{}) {
//	val = append([]interface{}{findCaller(0)}, val...)
//	err := glg.Error(val...)
//	if err != nil {
//		log.Println(err)
//		return
//	}
//}
func Println(val ...interface{}) {
	val = append([]interface{}{findCaller(0)}, val...)
	err := glg.Println(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func findCaller(skip int) string {
	file := ""
	line := 0
	for i := 0; i < 10; i++ {
		file, line = getCaller(skip + i)
		//log.Println(file,line)
		if !strings.HasPrefix(file, "logs.go") {
			break
		}
	}
	return fmt.Sprintf("(%s:%d)", file, line)
}
func getCaller(skip int) (string, int) {
	_, file, line, ok := runtime.Caller(skip)
	if !ok {
		return "", 0
	}
	n := 0
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			n++
			if n >= 2 {
				file = file[i+1:]
				break
			}
		}
	}
	return file, line
}
