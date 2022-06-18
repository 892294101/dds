package main

import (
	"context"
	"fmt"
	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/replication"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"os"
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
func InitHlog() *logrus.Logger {
	hlog := logrus.New()
	log := filepath.Join(GetHomeDirectory(), "logs", "extract.log")
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

	hlog.SetLevel(logrus.DebugLevel)
	hlog.SetReportCaller(true)
	hlog.SetFormatter(new(MyFormatter))
	hlog.Infof("Initialize log file: %s", log)

	return hlog
}

//根据执行文件路径获取程序的HOME路径
func GetHomeDirectory() (homedir string) {
	file, _ := exec.LookPath(os.Args[0])
	ExecFilePath, _ := filepath.Abs(file)

	os := runtime.GOOS
	switch os {
	case "windows":
		execfileslice := strings.Split(ExecFilePath, `\`)
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for i, v := range HomeDirectory {
			if v != "" {
				if i > 0 {
					homedir += `\` + v
				} else {
					homedir += v
				}
			}
		}
	case "linux":
		execfileslice := strings.Split(ExecFilePath, "/")
		HomeDirectory := execfileslice[:len(execfileslice)-2]
		for _, v := range HomeDirectory {
			if v != "" {
				homedir += `/` + v
			}
		}
	default:
		logrus.Error(fmt.Sprintf("Unknown operation type: %s", os))
	}

	if homedir == "" {
		logrus.Error(fmt.Sprintf("Get program home directory failed: %s", homedir))
	}
	return homedir
}

func main() {

	log := InitHlog()

	cfg := replication.BinlogSyncerConfig{
		ServerID: 100,
		Flavor:   "mysql",
		Host:     "127.0.0.1",
		Port:     3306,
		User:     "root",
		Password: "",
		Logger:   log,
	}

	syncer := replication.NewBinlogSyncer(cfg)

	streamer, err := syncer.StartSync(mysql.Position{"1", 1})
	if err != nil {
		log.Fatalf("%s", err)
	}

	for {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		ev, err := streamer.GetEvent(ctx)
		cancel()

		if err == context.DeadlineExceeded {
			// meet timeout
			continue
		}

		ev.Dump(os.Stdout)
	}
}
