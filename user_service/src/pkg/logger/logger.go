package logger

import "go.uber.org/zap"

var Log *zap.SugaredLogger

func Init() {
	logger := zap.Must(zap.NewDevelopment())
	Log = logger.Sugar()
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
