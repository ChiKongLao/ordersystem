package logger

import (
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/context"
	"time"
	"os"
	"path"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/sirupsen/logrus"
	"github.com/pkg/errors"
	"github.com/rifflock/lfshook"
)

// get a filename based on the date, file logs works that way the most times
// but these are just a sugar.
func todayFilename() string {
	today := time.Now().Format("2006-01-02")
	return today + ".txt"
}

func NewLogFile() *os.File {
	filename := todayFilename()
	// open an output file, this will append to the today's file if server restarted.
	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	return f
}

func GetRequestLogger() context.Handler {
	return logger.New(logger.Config{
		// Status displays status code
		Status: true,
		// IP displays request's remote address
		IP: true,
		// Method displays the http method
		Method: true,
		// Path displays the request path
		Path: true,

		//Columns: true,

		// if !empty then its contents derives from `ctx.Values().Get("logger_message")
		// will be added to the logs.
		MessageContextKey: "api_message",
	})
}

func ConfigLogger() {
	//ConfigLoggerWithParams(constant.ROOT,constant.LogFileName,1000000,20000)

	// You could set this to any `io.Writer` such as a file
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}
}

func ConfigLoggerWithParams(logPath string, logFileName string, maxAge time.Duration, rotationTime time.Duration) {
	baseLogPaht := path.Join(logPath, logFileName)
	writer, err := rotatelogs.New(
		baseLogPaht+".%Y%m%d%H%M",
		//rotatelogs.WithLinkName(baseLogPaht), // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		logrus.Errorf("config local file system logger error. %+v", errors.WithStack(err))
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer, // 为不同级别设置不同的输出目的
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, &logrus.TextFormatter{})
	logrus.AddHook(lfHook)
}
