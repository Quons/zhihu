package logging

import (
	"bytes"
	"fmt"
	"github.com/Quons/go-gin-example/pkg/file"
	"github.com/Quons/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin/json"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
	"path/filepath"
	"strconv"
	"time"
)

var logPath string
var dirName string

func Setup() {
	//获取执行目录
	var err error
	logPath, err = file.MkRdir("logs")
	if err != nil {
		logrus.Fatal("get log path error")
	}
	dirName, err = file.GetDirName()
	if err != nil {
		logrus.Fatal("get dirName error")
	}
	//设置日志级别
	logLevel, err := logrus.ParseLevel(setting.AppSetting.LogLevel)
	if err != nil {
		logrus.Fatal(err.Error())
	}
	logrus.SetLevel(logLevel)
	//打印行号，funcName
	logrus.SetReportCaller(true)
	//输出设置
	writer := GetLogrusWriter()
	//设置local file system hook
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &CodeFormatter{})
	//添加hook
	logrus.AddHook(lfsHook)
}

//定义formatter ,实现logrus formatter接口
type CodeFormatter struct{}

func (f *CodeFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	fileLineNum := ""
	if entry.Caller != nil {
		execPath, _ := file.GetExecPath()
		fileLineNum = string([]rune(entry.Caller.File)[len(execPath)+1:])
		fileLineNum = fmt.Sprintf("%s:%v ", fileLineNum, strconv.Itoa(entry.Caller.Line))
	}

	b.WriteString(entry.Time.Format("2006-01-02 15:04:05"))
	b.WriteString(" [")
	b.WriteString(entry.Level.String())
	b.WriteString("] ")
	b.WriteString(fileLineNum)

	if len(entry.Data) != 0 {
		b.WriteString("param:")
		data, _ := json.Marshal(entry.Data)
		b.WriteString(fmt.Sprintf("%+v ", string(data)))
	}

	b.WriteString("msg:")
	b.WriteString(entry.Message)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

func GetLogrusWriter() *rotatelogs.RotateLogs {
	logrusPath, err := file.MkRdir("logs/" + dirName)
	if err != nil {
		logrus.Fatal("get log path error")
	}
	writer, err := rotatelogs.New(
		filepath.Join(logrusPath, dirName+".%Y%m%d%H%M"),
		rotatelogs.WithLinkName(filepath.Join(logPath, dirName+".log")), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(10*time.Hour*24),                          // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24),                       // 日志切割时间间隔
	)
	if err != nil {
		logrus.Fatalf("config local file system logger error.%+v", errors.WithStack(err))
	}
	return writer
}

//获取gin日志writer
func GetGinLogWriter() *rotatelogs.RotateLogs {
	ginLogPath, err := file.MkRdir("logs/gin")
	if err != nil {
		logrus.Fatal("get log path error")
	}
	writer, err := rotatelogs.New(
		filepath.Join(ginLogPath, "gin.%Y%m%d%H%M"),
		rotatelogs.WithLinkName(filepath.Join(logPath, "gin.log")), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(10*time.Hour*24),                     // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24),                  // 日志切割时间间隔
	)
	if err != nil {
		logrus.Fatalf("config local file system logger error.%+v", errors.WithStack(err))
	}
	return writer
}

//获取gorm日志writer
func GetGormLogWriter() *rotatelogs.RotateLogs {
	writer, err := rotatelogs.New(
		filepath.Join(logPath, "gorm.%Y%m%d%H%M"),
		rotatelogs.WithLinkName(filepath.Join(logPath, "gorm.log")), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(10*time.Hour*24),                      // 文件最大保存时间
		rotatelogs.WithRotationTime(time.Hour*24),                   // 日志切割时间间隔
	)
	if err != nil {
		logrus.Fatalf("config local file system logger error.%+v", errors.WithStack(err))
	}
	return writer
}
