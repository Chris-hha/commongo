package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"os"
	"time"
)

const (
	logDir = "logs"
	logName      = "saber_cfg_render.log"
	errLogName   = "saber_cfg_render_err.log"
	logLevel     = "info"
	maxRemainCnt = 7
)

var logLevels = map[string]logrus.Level{
	"panic": logrus.PanicLevel,
	"fatal": logrus.FatalLevel,
	"error": logrus.ErrorLevel,
	"warn":  logrus.WarnLevel,
	"info":  logrus.InfoLevel,
	"debug": logrus.DebugLevel,
	"trace": logrus.TraceLevel,
}

type ErrorLogHook struct {
	Writer *os.File
}

func (elh *ErrorLogHook) Fire(entry *logrus.Entry) error {
	line, err := entry.String()
	if err != nil {
		logrus.Errorf("Unable to read entry, %v", err)
		return err
	}
	_, err = elh.Writer.Write([]byte(line))
	if err != nil {
		logrus.Errorf("write line error, line: %v, err: %v", line, err)
	}
	return nil
}

func (elh *ErrorLogHook) Levels() []logrus.Level {
	return []logrus.Level{logLevels["error"], logLevels["fatal"], logLevels["panic"]}
}

func NewElh() *ErrorLogHook {
	// error log单独存储
	ef, err := os.OpenFile(fmt.Sprintf("%s/%s", logDir, errLogName), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to create error logfile" + errLogName)
		panic(err)
	}
	// defer ef.Close()
	return &ErrorLogHook{
		Writer: ef,
	}
}

func newLfsHook() logrus.Hook {
	writer, err := rotatelogs.New(
		logDir + "/" + logName + "-%Y%m%d%H",
		// WithLinkName为最新的日志建立软连接,以方便随着找到当前日志文件
		rotatelogs.WithLinkName(logDir + "/" + logName),
		// WithRotationTime设置日志分割的时间,这里设置为一小时分割一次
		rotatelogs.WithRotationTime(time.Hour*24),
		// WithMaxAge和WithRotationCount二者只能设置一个,
		// WithMaxAge设置文件清理前的最长保存时间,
		// WithRotationCount设置文件清理前最多保存的个数.
		// rotatelogs.WithMaxAge(time.Hour*24),
		rotatelogs.WithRotationCount(maxRemainCnt),
	)
	if err != nil {
		logrus.Errorf("Init newLfsHook writer for logger error: %v", err)
	}
	// 判断输入的日志等级是否存在，不存在则给一个默认值
	if level, ok := logLevels[logLevel]; ok {
		logrus.SetLevel(level)
	} else {
		logrus.SetLevel(logLevels[logLevel])
	}

	// 使用了lfshook软件包创建了一个新的日志钩子，该钩子将日志记录到指定的日志文件中。
	// lfshook.WriterMap指定了每个日志级别所使用的写入器（writer）。
	// 在这个函数中，所有的日志级别都使用同一个写入器writer。
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{
		ForceColors:     true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	return lfsHook
}

func init() {
	// 创建一个Logrus的实例
	logger := logrus.New()

	// 输出行号
	logger.SetReportCaller(true)

	err := os.MkdirAll(logDir, os.ModePerm)
	if err != nil {
		logrus.Error("make logDir error")
		return
	}

	// 错误日志单独存，7天分割
	elh := NewElh()
	lfsHook := newLfsHook()
	logger.AddHook(elh)
	logger.AddHook(lfsHook)
}

