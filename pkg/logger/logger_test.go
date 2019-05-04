package logger

import (
	"testing"

	"go.uber.org/zap"

	"github.com/nikitaksv/bodis/internal/integration/yandexdisk"

	"github.com/nikitaksv/bodis/pkg/config"

	"github.com/nikitaksv/bodis/pkg/storage"
)

func TestStorageError(t *testing.T) {
	type args struct {
		err   storage.Error
		sugar *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{"without params", args{storage.NewBaseError(yandexdisk.StorageKey, "INVALIDTOKEN", "invalid token", nil), Sugar()}},
		{"with params", args{storage.NewBaseError(yandexdisk.StorageKey, "INVALIDTOKEN", "invalid token",
			map[string]interface{}{
				"method": "GetResourceUploadLink()",
				"arguments": map[string]interface{}{
					"path": "/tests/foo.png",
					"data": map[string]int{
						"len": 122,
						"cap": 244,
					},
				},
			},
		), Sugar()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			StorageError(tt.args.err, tt.args.sugar)
		})
	}
}

func TestConfigFileError(t *testing.T) {
	type args struct {
		err   *config.FileError
		sugar *zap.SugaredLogger
	}
	tests := []struct {
		name string
		args args
	}{
		{"1", args{&config.FileError{FilePath: "/bodis/cmd/config.toml", Message: "invalid token", Line: 23, Key: "token", Value: "aaaaaa11111222222333333", Table: "storage"}, Sugar()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ConfigFileError(tt.args.err, tt.args.sugar)
			sync(tt.args.sugar)
		})
	}
}
