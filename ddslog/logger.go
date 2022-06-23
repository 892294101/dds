package ddslog

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"myGithubLib/dds/extract/mysql/utils"
	"os"
	"path/filepath"
	"strings"
	"time"
)

//自义定日志结构
type MyFormatter struct{}

func (s *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := time.Now().Local().Format("2006-01-02 15:04:05.00")
	var reason interface{}
	if v, ok := entry.Data["err"]; ok {
		reason = v
	} else {
		reason = nil
	}
	var msg string
	fName := entry.Caller.Function[strings.LastIndex(entry.Caller.Function, ".")+1:]
	if reason == nil {
		msg = fmt.Sprintf("%s %s %s %v [M] %s\n", timestamp, strings.ToUpper(entry.Level.String())[:1], fName, entry.Caller.Line, entry.Message)
	} else {
		msg = fmt.Sprintf("%s %s %s %d [M] %s %v\n", timestamp, strings.ToUpper(entry.Level.String())[:1], fName, entry.Caller.Line, entry.Message, reason)
	}
	return []byte(msg), nil
}

//初始化日志输出
//定义为同事输出日志内容到标准输出和和日志文件
func InitDDSlog() (*logrus.Logger, error) {
	ddslog := logrus.New()
	dir, err := utils.GetHomeDirectory()
	if err != nil {
		return nil, err
	}
	log := filepath.Join(*dir, "logs", "extract.log")
	logfile, err := os.OpenFile(log, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	//writers := []io.Writer{logfile, os.Stdout}
	writers := []io.Writer{logfile}
	fileAndStdoutWriter := io.MultiWriter(writers...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Log file open failed: %s", err)
		os.Exit(1)
	} else {
		logrus.SetOutput(fileAndStdoutWriter)
	}

	ddslog.SetLevel(logrus.DebugLevel)
	ddslog.SetReportCaller(true)
	ddslog.SetFormatter(new(MyFormatter))
	ddslog.Infof("Initialize log file: %s", log)

	return ddslog, nil
}
