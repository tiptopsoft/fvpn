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

func Log() *zap.SugaredLogger {
	return zapLogger
}
