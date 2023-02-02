package log

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func init() {
	logger, _ := zap.NewProduction()
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}(logger) // flushes buffer, if any
	Logger = logger.Sugar()

}
