package logger

import (
	"go.uber.org/zap"

	"github.com/nikitaksv/bodis/pkg/errors"
)

func sync(zap *zap.Logger) {
	err := zap.Sync()
	if err != nil {
		panic(err)
	}
}

func Sugar() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer sync(logger)
	return logger.Sugar()
}

func Error(err errors.Error) {
	sugar := Sugar()
	// TODO реализовать неймспейс
	sugar.With(zap.Namespace())
	sugar.Errorw(err.ErrorID(), zap.Any("description", err.Description()), zap.Any("params", err.Params()))
}
