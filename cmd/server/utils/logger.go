package utils

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/app/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Logger *zap.Logger

func NewLogger(appEnv string) *zap.Logger {
	switch appEnv {
	case constants.ProductionMode:
		core := zapcore.NewTee(
			DevelopmentLogger(),
			ProductionLogger(),
		)
		return zap.New(core)

	case constants.DevelopmentMode:
		core := zapcore.NewTee(
			DevelopmentLogger(),
			ProductionLogger(),
		)
		return zap.New(core)
	default:
		core := zapcore.NewTee(
			DevelopmentLogger(),
			ProductionLogger(),
		)
		return zap.New(core)
	}
}

func DevelopmentLogger() zapcore.Core {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   constants.LogFilePath,
		MaxSize:    2, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})
	level := zap.NewAtomicLevelAt(zap.DebugLevel)
	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	developmentCfg.TimeKey = constants.TimeKey
	developmentCfg.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimestampFormat)
	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	return zapcore.NewCore(consoleEncoder, file, level)
}

func ProductionLogger() zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)
	level := zap.NewAtomicLevelAt(zap.ErrorLevel)
	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = constants.TimeKey
	productionCfg.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimestampFormat)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)
	return zapcore.NewCore(fileEncoder, stdout, level)
}
