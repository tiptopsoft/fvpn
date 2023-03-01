package log

import (
	"fmt"
	"go.uber.org/zap"
)

var zapLogger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}(logger) // flushes buffer, if any
	zapLogger = logger.Sugar()
}

func Errorf(templates string, args ...interface{}) {
	zapLogger.Errorf(templates, args)
}

func Error(args ...interface{}) {
	zapLogger.Error(args)
}

func Infof(templates string, args ...interface{}) {
	zapLogger.Infof(templates, args)
}

func Info(args ...interface{}) {
	zapLogger.Info(args)
}

func Debugf(templates string, args ...interface{}) {
	zapLogger.Debugf(templates, args)
}

func Warnf(templates string, args ...interface{}) {
	zapLogger.Warnf(templates, args)
}
