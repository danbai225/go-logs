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
)

var logsDir = "./logs"
var logsDirC = make(chan string, 1)
var httpLog, infoLog, errLog, debugLog, warnLog *os.File
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
}

func init() {
	loadFiles()
	glg.Get().SetMode(glg.BOTH)
	glg.Get().SetLineTraceMode(glg.TraceLineNone)
	glg.Get().SetTimeLocation(time.Now().Location())
	go splitLogByDay()
}

var logsFiles []*os.File
var ini = true

func timeRemaining() time.Duration {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02") + " 00:00:00"
	tomorrowTime, _ := time.ParseInLocation("2006-01-02 15:04:05", tomorrow, time.Local)
	return time.Until(tomorrowTime)
}
func splitLog(fromTime bool) {
	cuttingOff()
	closeFiles()
	//上一天日期
	previousDate := time.Now().Add(time.Minute * -1).Format("2006-01-02")
	if (writeLogs & INFO) != 0 {
		wt := io.Writer(infoLog)
		if fromTime {
			wt = io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-info.log", logsDir, previousDate)), infoLog)
		}
		glg.Get().SetLevelWriter(glg.INFO, wt)
	}
	if (writeLogs & ERR) != 0 {
		wt := io.Writer(errLog)
		if fromTime {
			wt = io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-err.log", logsDir, previousDate)), errLog)
		}
		glg.Get().SetLevelWriter(glg.ERR, wt)
	}
	if (writeLogs & WARN) != 0 {
		wt := io.Writer(warnLog)
		if fromTime {
			wt = io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-warn.log", logsDir, previousDate)), warnLog)
		}
		glg.Get().SetLevelWriter(glg.WARN, wt)
	}
	if (writeLogs & DEBUG) != 0 {
		wt := io.Writer(debugLog)
		if fromTime {
			wt = io.MultiWriter(getFileWriter(fmt.Sprintf("%s/%s-debug.log", logsDir, previousDate)), debugLog)
		}
		glg.Get().SetLevelWriter(glg.DEBG, wt)
	}
	if fromTime {
		go clearingRedundantLogs()
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
func DebugF(format string, a ...interface{}) {
	DebugN(0, fmt.Sprintf(format, a...))
}
func InfoF(format string, a ...interface{}) {
	InfoN(0, fmt.Sprintf(format, a...))
}
func WarnF(format string, a ...interface{}) {
	WarnN(0, fmt.Sprintf(format, a...))
}
func ErrF(format string, a ...interface{}) {
	ErrN(0, fmt.Sprintf(format, a...))
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
func DebugN(n int, val ...interface{}) {
	val = append([]interface{}{findCaller(n + 3)}, val...)
	err := glg.Debug(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Debug(val ...interface{}) {
	DebugN(0, val...)
}
func InfoN(n int, val ...interface{}) {
	val = append([]interface{}{findCaller(n + 3)}, val...)
	err := glg.Info(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Info(val ...interface{}) {
	InfoN(0, val...)
}
func WarnN(n int, val ...interface{}) {
	val = append([]interface{}{findCaller(n + 3)}, val...)
	err := glg.Warn(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Warn(val ...interface{}) {
	WarnN(0, val...)
}
func ErrN(n int, val ...interface{}) {
	val = append([]interface{}{findCaller(n + 3)}, val...)
	err := glg.Error(val...)
	if err != nil {
		log.Println(err)
		return
	}
}
func Err(val ...interface{}) {
	ErrN(0, val...)
}
func SetLevel(l int) {
	switch l {
	case DEBUG:
		glg.Get().SetLevel(glg.DEBG)
	case INFO:
		glg.Get().SetLevel(glg.INFO)
	case WARN:
		glg.Get().SetLevel(glg.WARN)
	case ERR:
		glg.Get().SetLevel(glg.ERR)
	default:
		glg.Get().SetLevel(glg.INFO)
	}
}
