package log

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	Logger = logger.Sugar()

}
