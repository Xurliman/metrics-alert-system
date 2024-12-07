package log

import (
	"fmt"
	"github.com/Xurliman/metrics-alert-system/cmd/agent/app/constants"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logger   *zap.Logger
	basePath string
)

func init() {
	_, b, _, _ := runtime.Caller(0)
	basePath = filepath.Dir(filepath.Dir(filepath.Dir(b)))
}

func InitLogger(environment string, logPath string) {
	logger = initLogger(environment, logPath)
}

func Info(msg string, fields ...zap.Field) {
	getLoggerWithSkip(1).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	getLoggerWithSkip(1).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	getLoggerWithSkip(1).Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	getLoggerWithSkip(1).Fatal(msg, fields...)
}

func getLoggerWithSkip(skip int) *zap.Logger {
	return logger.WithOptions(zap.AddCallerSkip(skip))
}

func initLogger(environment string, logPath string) *zap.Logger {
	var core zapcore.Core

	consoleEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:          "time",
		LevelKey:         "level",
		NameKey:          "logger",
		CallerKey:        "caller",
		MessageKey:       "message",
		StacktraceKey:    "stacktrace",
		LineEnding:       zapcore.DefaultLineEnding,
		EncodeLevel:      zapcore.CapitalColorLevelEncoder,
		EncodeTime:       zapcore.ISO8601TimeEncoder,
		EncodeDuration:   zapcore.StringDurationEncoder,
		EncodeCaller:     shortCallerEncoder,
		ConsoleSeparator: " | ",
	})

	productionEncoderConfig := zap.NewProductionEncoderConfig()
	productionEncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(constants.TimestampFormat)
	productionEncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	productionEncoderConfig.TimeKey = "time"
	productionEncoderConfig.MessageKey = "message"
	productionEncoderConfig.ConsoleSeparator = " | "
	productionEncoderConfig.EncodeCaller = shortCallerEncoder
	fileEncoder := zapcore.NewJSONEncoder(productionEncoderConfig)

	logFile := &lumberjack.Logger{
		Filename:   filepath.Join(basePath, logPath),
		MaxSize:    10, // MB
		MaxBackups: 5,  // old log files to keep
		MaxAge:     30, // days
		Compress:   true,
	}

	logLevel := zapcore.DebugLevel
	if environment == "production" {
		logLevel = zapcore.InfoLevel
	}
	core = zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), logLevel),
		zapcore.NewCore(fileEncoder, zapcore.AddSync(logFile), logLevel),
	)

	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
}

func shortCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%s:%d", filepath.Base(caller.File), caller.Line))
}
