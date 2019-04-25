package yandexDisk

import (
	"context"
	"net/http"

	"github.com/nikitaksv/yandex-disk-sdk-go"

	"github.com/nikitaksv/bodis/pkg/storage"
)

type yandexDisk struct {
	client yadisk.YaDisk
}

func NewYandexDisk(ctx context.Context, client http.Client, token string) *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(ctx, &client, &yadisk.Token{AccessToken: token})
	if err != nil {
		panic(err)
	}

	return &yandexDisk{client: yadClient}
}

func (yd *yandexDisk) Info() (storage.Info, error) {
	disk, err := yd.client.GetDisk([]string{"total_space", "max_file_size", "used_space"})
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, e
		}
		panic(err)
	}
	return newYandexDiskInfo(uint64(disk.TotalSpace), uint64(disk.UsedSpace), uint64(disk.MaxFileSize)), nil
}

func (yd *yandexDisk) GetResourceInfo(path string) (storage.ResourceInfo, error) {
	res, err := yd.client.GetResource(path, nil, 10, 0, false, "", "size")
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, e
		}
		panic(err)
	}
	res.Embedded.Items

}

func (yd *yandexDisk) ReadResource(path string) ([]byte, error) {
	panic("implement me")
}
func (yd *yandexDisk) WriteResource(path string, data []byte) error {
	panic("implement me")
}

func (yd *yandexDisk) DeleteResource(path string) error {
	panic("implement me")
}

func (yd *yandexDisk) StateResource(path string) {
	panic("implement me")
}
