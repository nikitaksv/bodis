package integration

import (
	"context"
	"net/http"

	"github.com/nikitaksv/bodis/internal/integration/yandexdisk"

	"github.com/nikitaksv/bodis/pkg/storage"
)

type Settings map[string]interface{}
type Constructor struct {
	IntegrationKey storage.IntegrationKey
	AuthData       storage.AuthData
	Settings       Settings
}

func GetStorage(constructor Constructor) storage.Storage {
	switch constructor.IntegrationKey {
	case yandexdisk.StorageKey:
		return createYandexDiskIntegration(constructor.Settings, constructor.AuthData)
	}
	return nil
}

func createYandexDiskIntegration(settings Settings, authData storage.AuthData) storage.Storage {
	return yandexdisk.New(settings["context"].(context.Context), settings["client"].(*http.Client), authData["token"])
}
