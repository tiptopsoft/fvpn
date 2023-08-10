// Copyright 2023 Tiptopsoft, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package log

import (
	"github.com/tiptopsoft/fvpn/pkg/util"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

var zaplogger *zap.SugaredLogger

func init() {
	var logger *zap.Logger
	config, err := util.InitConfig()
	if err != nil {
		panic(err)
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	// 设置日志记录中时间的格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// 日志Encoder 还是JSONEncoder，把日志行格式化成JSON格式的
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	var core zapcore.Core
	if config.NodeCfg.Log.EnableDebug {
		core = zapcore.NewTee(
			// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是Debug 及以上，生产环境记得改成Info
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		)
	} else {
		core = zapcore.NewTee(
			// 同时向控制台和文件写日志， 生产环境记得把控制台写入去掉，日志记录的基本是Debug 及以上，生产环境记得改成Info
			zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.InfoLevel),
		)
	}

	logger = zap.New(core)
	defer logger.Sync()
	zaplogger = logger.Sugar()
}

func Log() *zap.SugaredLogger {
	return zaplogger
}
