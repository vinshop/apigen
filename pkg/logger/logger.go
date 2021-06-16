package logger

import "go.uber.org/zap"

var logger *zap.SugaredLogger

func init() {
	cfg := zap.NewDevelopmentConfig()
	cfg.OutputPaths = []string{
		"apigen.log",
	}
	l, err := cfg.Build()
	if err != nil {
		zap.S().Fatalw("Error when build log", "error", err)
	}
	logger = l.Sugar()
	Info("Starting generate")
}

func Error(args ...interface{}) {
	logger.Error(args...)
}

func Errorw(msg string, args ...interface{}) {
	logger.Errorw(msg, args...)
}

func Info(args ...interface{}) {
	logger.Info(args...)
}

func Infow(msg string, args ...interface{}) {
	logger.Infow(msg, args...)
}

func Warnw(msg string, args ...interface{}) {
	logger.Warnw(msg, args...)
}

func Fatalw(msg string, args ...interface{}) {
	logger.Fatalw(msg, args...)
}
