package logger

import (
	"go.uber.org/zap"

	"github.com/nikitaksv/bodis/pkg/storage"
)

func sync(zap *zap.Logger) {
	err := zap.Sync()
	if err != nil {
		panic(err)
	}
}

func Sugar() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer sync(logger) // flushes buffer, if any
	return logger.Sugar()
}

func Error(err storage.Error) {
	sugar := Sugar()
	sugar.With(zap.Namespace())
	sugar.Errorw(err.ErrorID(), zap.Any("description", err.Description()), zap.Any("params", err.Params()))
}
