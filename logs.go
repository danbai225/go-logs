package go_logs

import (
	"fmt"
	"github.com/kpango/glg"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	DEBUG = 1
	INFO  = 2
	WARN  = 4
	ERR   = 8
)

var logsDir = "logs"

// GinLog gin.DefaultWriter = io.MultiWriter(inits.GinLog, os.Stdout)
var GinLog, infoLog, errLog, debugLog, warnLog, stdErrLog *os.File
var writeLogs = byte(INFO | ERR)
var StderrFile = false

func SetLogsDir(dir string) {
	logsDir = dir
	loadFiles()
}
func SetWriteLogs(logs byte) {
	writeLogs = logs
}
func loadFiles() {
	GinLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "gin.log"), 0777)
	infoLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log"), 0777)
	errLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "error.log"), 0777)
	debugLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log"), 0777)
	warnLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log"), 0777)
	stdErrLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "stdErr.log"), 0777)

}

func init() {
	loadFiles()
	glg.Get().SetMode(glg.BOTH)
	glg.Get().SetLineTraceMode(glg.TraceLineNone)
	go splitLogByDay()
	go rewriteStderrFile()
}
func splitLogByDay() {
	timeDay := "2006-01-02"
	cDir := logsDir
	for {
		NowTimeDay := time.Now().Format("2006-01-02")
		if NowTimeDay > timeDay || cDir != logsDir {
			cuttingOff()
			if (writeLogs & INFO) != 0 {
				glg.Get().SetLevelWriter(glg.INFO, io.MultiWriter(glg.FileWriter(fmt.Sprintf("%s/%s-info.log", logsDir, NowTimeDay), 0777), infoLog))
			}
			if (writeLogs & ERR) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("%s/%s-err.log", logsDir, NowTimeDay), 0777), errLog))
			}
			if (writeLogs & WARN) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("%s/%s-warn.log", logsDir, NowTimeDay), 0777), warnLog))
			}
			if (writeLogs & DEBUG) != 0 {
				glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(glg.FileWriter(fmt.Sprintf("%s/%s-debug.log", logsDir, NowTimeDay), 0777), debugLog))
			}
		}
		time.Sleep(time.Second)
	}
}
func Debug(val ...interface{}) {
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Debug(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Info(val ...interface{}) {
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Info(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Warn(val ...interface{}) {
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Warn(val...)
	if err != nil {
		log.Println(err)
		return
	}
}

func Err(val ...interface{}) {
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Error(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Println(val ...interface{}) {
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Println(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func findCaller(skip int) string {
	var fl string
	_, file, line, ok := runtime.Caller(skip)
	switch {
	case !ok:
		fl = "???:0"
	case strings.HasPrefix(file, runtime.GOROOT()+"/src"):
		fl = "https://github.com/golang/go/blob/" + runtime.Version() + strings.TrimPrefix(file, runtime.GOROOT()) + "#L" + strconv.Itoa(line)
	case strings.Contains(file, "go/pkg/mod/"):
		fl = "https:/"
		for _, path := range strings.Split(strings.SplitN(file, "go/pkg/mod/", 2)[1], "/") {
			if strings.Contains(path, "@") {
				sv := strings.SplitN(path, "@", 2)
				if strings.Count(sv[1], "-") > 2 {
					path = sv[0] + "/blob/master"
				} else {
					path = sv[0] + "/blob/" + sv[1]
				}
			}
			fl += "/" + path
		}
		fl += "#L" + strconv.Itoa(line)
	case strings.Contains(file, "go/src"):
		fl = "https:/"
		cnt := 0
		for _, path := range strings.Split(strings.SplitN(file, "go/src/", 2)[1], "/") {
			if cnt == 3 {
				path = "blob/master/" + path
			}
			fl += "/" + path
			cnt++
		}
		fl += "#L" + strconv.Itoa(line)
	default:
		for i := len(file) - 1; i > 0; i-- {
			if file[i] == '/' {
				file = file[i+1:]
				break
			}
		}
		fl = file + ":" + strconv.Itoa(line)
	}
	return fl
}

//func getCaller(skip int) (string, int) {
//	_, file, line, ok := runtime.Caller(skip)
//	if !ok {
//		return "", 0
//	}
//	n := 0
//	for i := len(file) - 1; i > 0; i-- {
//		if file[i] == '/' {
//			n++
//			if n >= 2 {
//				file = file[i+1:]
//				break
//			}
//		}
//	}
//	return file, line
//}
