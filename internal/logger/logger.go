package logger

import (
	"fmt"
	"tasktracker/configs"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type (
	Logger struct {
		L *zap.Logger
	}
)

const (
	// filepath defines when logs reside
	filepath = `../../logs/app.log`
)

// Init creates and initializes logger and returns it and an error if any
func Init(cfg *configs.Config) (*Logger, error) {
	logLevel, err := getLogLevel(cfg)
	if err != nil {
		return nil, fmt.Errorf("defining log level: %w", err)
	}
	encoderCfg := getEncoderConfig(logLevel)

	subLogger := configurateSublogger(cfg)
	writeSyncer := zapcore.AddSync(subLogger)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writeSyncer,
		logLevel,
	)

	return &Logger{L: zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.DebugLevel))}, nil
}

// configurateSublogger builds and returns sublogger based on lumberjack.Logger
func configurateSublogger(cfg *configs.Config) *lumberjack.Logger {
	logger := &lumberjack.Logger{
		Filename:   filepath,
		MaxSize:    cfg.MaxFileSize,
		MaxAge:     cfg.MaxAge,
		MaxBackups: cfg.MaxBackups,
		LocalTime:  cfg.IsLocalTime,
		Compress:   cfg.IsCompressed,
	}

	return logger
}

// getLogLevel checks the value of cfg.Status and returns a corresponding zapcore.Level and an error if any
func getLogLevel(cfg *configs.Config) (zapcore.Level, error) {
	const (
		debugStatus = "debug"
		devStatus   = "dev"
		prodStatus  = "prod"
	)

	switch cfg.Environment.Status {
	case debugStatus:
		return zapcore.DebugLevel, nil

	case devStatus, prodStatus:
		return zapcore.InfoLevel, nil

	default:
		return zapcore.InvalidLevel, fmt.Errorf("defining log level: unrecognized env variable")
	}
}

// getEncoderConfig builds a zapcore.EncoderConfig based on logLevel and returns it
func getEncoderConfig(logLevel zapcore.Level) zapcore.EncoderConfig {
	var (
		encoderCfg zapcore.EncoderConfig
	)

	switch logLevel {
	case zapcore.DebugLevel:
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	case zapcore.InfoLevel:
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	return encoderCfg
}
