package integration

import (
	"github.com/nikitaksv/bodis/internal/integration/yandexDisk"
	"github.com/nikitaksv/bodis/pkg/storage"
)

type IntegrationKey string

type Constructor struct {
	IntegrationKey IntegrationKey
	AuthData       map[string]string
	Settings       interface{}
}

func GetStorage(constructor Constructor) storage.Storage {
	switch constructor.IntegrationKey {
	case yandexDisk.IntegrationKey:
		return yandexDisk.NewYandexDisk()
	}
}
