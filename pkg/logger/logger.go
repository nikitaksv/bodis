package logger

import (
	"go.uber.org/zap"

	"github.com/nikitaksv/bodis/pkg/config"

	"github.com/nikitaksv/bodis/pkg/storage"
)

func sync(zap *zap.SugaredLogger) {
	err := zap.Sync()
	if err != nil {
		panic(err)
	}
}

func Sugar() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	return logger.Sugar()
}

func StorageError(err storage.Error, sugar *zap.SugaredLogger) {
	sugar.Errorw(err.ErrorID(),
		zap.String("storageKey", string(err.StorageKey())),
		zap.Any("description", err.Description()),
		zap.Any("params", err.Params()),
	)
}

func ConfigFileError(err *config.FileError, sugar *zap.SugaredLogger) {
	sugar.Errorw(err.Message,
		zap.String("line", string(err.Line)),
		zap.String("filepath", err.FilePath),
		zap.String("key", err.Key),
		zap.String("value", err.Value),
	)
}
