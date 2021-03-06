package go_logs

import (
	"encoding/json"
	"fmt"
	"github.com/kpango/glg"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
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
	STD   = 16
)

var logsDir = "./logs"
var logsDirC = make(chan string, 1)
var stdLog, httpLog, infoLog, errLog, debugLog, warnLog *os.File
var writeLogs = byte(INFO | ERR)
var saveDay = 30
var redirectStdLog = false
var after = time.After(timeRemaining())

// SetLogsDir 切换目录可能会会损失部分日志
func SetLogsDir(dir string) {
	if dir != logsDir {
		logsDirC <- dir
	}
	logsDir = dir
	loadFiles()
}
func SetWriteLogs(logs byte) {
	writeLogs = logs
	if redirectStdLog {
		writeLogs = writeLogs | STD
	}
	loadFiles()
}
func GetHttpWriter() *os.File {
	if httpLog == nil {
		httpLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "http.log"), 0777)
	}
	return httpLog
}
func SetSaveDay(day int) {
	saveDay = day
}
func SetRedirectStdLog() {
	redirectStdLog = true
	SetWriteLogs(writeLogs)
}
func loadFiles() {
	if (writeLogs & INFO) != 0 {
		infoLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "info.log"), 0777)
		glg.Get().SetLevelWriter(glg.INFO, infoLog)
	} else if infoLog != nil {
		_ = infoLog.Close()
	}
	if (writeLogs & ERR) != 0 {
		errLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "err.log"), 0777)
		glg.Get().SetLevelWriter(glg.ERR, errLog)
	} else if errLog != nil {
		_ = errLog.Close()
	}
	if (writeLogs & DEBUG) != 0 {
		debugLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "debug.log"), 0777)
		glg.Get().SetLevelWriter(glg.DEBG, debugLog)
	} else if debugLog != nil {
		_ = debugLog.Close()
	}
	if (writeLogs & WARN) != 0 {
		warnLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "warn.log"), 0777)
		glg.Get().SetLevelWriter(glg.WARN, warnLog)
	} else if warnLog != nil {
		_ = warnLog.Close()
	}
	if (writeLogs & STD) != 0 {
		stdLog = glg.FileWriter(fmt.Sprintf("%s%c%s", logsDir, os.PathSeparator, "std.log"), 0777)
		err := redirectStdout(stdLog)
		if err != nil {
			println(err.Error())
		}
		err = redirectStderr(stdLog)
		if err != nil {
			println(err.Error())
		}
	} else if stdLog != nil {
		_ = stdLog.Close()
	}
}

func init() {
	loadFiles()
	glg.Get().SetMode(glg.BOTH)
	glg.Get().SetLineTraceMode(glg.TraceLineNone)
	go splitLogByDay()
}

var logsFiles []*os.File
var ini = true

func timeRemaining() time.Duration {
	todayLast := time.Now().Format("2006-01-02") + " 23:59:59"
	todayLastTime, _ := time.ParseInLocation("2006-01-02 15:04:05", todayLast, time.Local)
	return time.Duration(todayLastTime.Unix()-time.Now().Local().Unix()) * time.Second
}
func splitLog(fromTime bool) {
	NowTimeDay := time.Now().Format("2006-01-02")
	cuttingOff()
	closeFiles()
	if (writeLogs & INFO) != 0 {
		glg.Get().SetLevelWriter(glg.INFO, io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-info.log", logsDir, NowTimeDay)), infoLog))
	}
	if (writeLogs & ERR) != 0 {
		glg.Get().SetLevelWriter(glg.ERR, io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-err.log", logsDir, NowTimeDay)), errLog))
	}
	if (writeLogs & WARN) != 0 {
		glg.Get().SetLevelWriter(glg.WARN, io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-warn.log", logsDir, NowTimeDay)), warnLog))
	}
	if (writeLogs & DEBUG) != 0 {
		glg.Get().SetLevelWriter(glg.DEBG, io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-debug.log", logsDir, NowTimeDay)), debugLog))
	}
	go clearingRedundantLogs()
	if fromTime {
		after = time.After(timeRemaining())
	}
}
func splitLogByDay() {
	for {
		select {
		case <-after:
			splitLog(true)
		case <-logsDirC:
			splitLog(false)
		}
	}
}
func getLogsName() []string {
	files := make([]string, 0)
	root := logsDir
	_ = filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	return files
}
func clearingRedundantLogs() {
	logs := getLogsName()
	//今日日期
	now := time.Now()
	date := now.Add(-(time.Duration(now.Hour())*time.Hour + time.Duration(now.Minute())*time.Minute + time.Duration(now.Second())*time.Second))
	inDate := make(map[string]struct{})

	for i := 0; i <= saveDay; i++ {
		d := date.AddDate(0, 0, -i)
		format := d.Format("2006-01-02")
		inDate[format] = struct{}{}
	}
	for _, logPath := range logs {
		base := path.Base(logPath)
		if len(base) > 10 {
			logDate := base[:10]
			if _, has := inDate[logDate]; !has {
				_ = os.Remove(logPath)
			}
		}
	}
}
func getFileWriter(path string) *os.File {
	if logsFiles == nil {
		logsFiles = make([]*os.File, 0)
	}
	writer := glg.FileWriter(path, 0666)
	logsFiles = append(logsFiles, writer)
	return writer
}
func closeFiles() {
	if logsFiles != nil {
		for _, file := range logsFiles {
			if file != nil {
				_ = file.Close()
			}
		}
	}
	logsFiles = make([]*os.File, 0)
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
func PrintJson(val ...interface{}) {
	var arr []string
	for i := range val {
		marshal, err := json.Marshal(val[i])
		if err != nil {
			marshal = []byte(err.Error())
		}
		arr = append(arr, string(marshal))
	}
	val = append([]interface{}{findCaller(2)}, val...)
	err := glg.Info(strings.Join(arr, ","))
	if err != nil {
		log.Println(err)
		return
	}
}
func InfoF(format string, a ...interface{}) {
	Info(fmt.Sprintf(format, a...))
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
func ErrF(format string, a ...interface{}) {
	Err(fmt.Sprintf(format, a...))
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
		for _, lPath := range strings.Split(strings.SplitN(file, "go/pkg/mod/", 2)[1], "/") {
			if strings.Contains(lPath, "@") {
				sv := strings.SplitN(lPath, "@", 2)
				if strings.Count(sv[1], "-") > 2 {
					lPath = sv[0] + "/blob/master"
				} else {
					lPath = sv[0] + "/blob/" + sv[1]
				}
			}
			fl += "/" + lPath
		}
		fl += "#L" + strconv.Itoa(line)
	case strings.Contains(file, "go/src"):
		fl = "https:/"
		cnt := 0
		for _, lPath := range strings.Split(strings.SplitN(file, "go/src/", 2)[1], "/") {
			if cnt == 3 {
				lPath = "blob/master/" + lPath
			}
			fl += "/" + lPath
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
