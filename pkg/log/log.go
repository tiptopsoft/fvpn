package log

import (
	"go.uber.org/zap"
)

var zapLogger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	zapLogger = logger.Sugar()
}

func Log() *zap.SugaredLogger {
	return zapLogger
}
