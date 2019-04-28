package yandexDisk

import (
	"context"
	"net/http"

	"github.com/nikitaksv/yandex-disk-sdk-go"

	"github.com/nikitaksv/bodis/pkg/logger"

	"github.com/nikitaksv/bodis/pkg/storage"
)

type yandexDisk struct {
	client yadisk.YaDisk
}

func NewYandexDisk(ctx context.Context, client *http.Client, token string) *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(ctx, client, &yadisk.Token{AccessToken: token})
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
	return yd.getResourceInfo(path)
}

func (yd *yandexDisk) getResourceInfo(path string) (*resourceInfo, error) {
	sugar := logger.Sugar()
	res, err := yd.client.GetResource(path, nil, 100, 0, false, "", "size")
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {

			return nil, e
		}
		panic(err)
	}

	ri := newResourceInfo(yd, res.Path, res.Name, res.Type, res.Md5, res.Created, res.Modified, uint64(res.Size), nil, nil)
	resources := make([]resourceInfo, len(res.Embedded.Items))
	for i, v := range res.Embedded.Items {
		resources[i] = *newResourceInfo(yd, v.Path, v.Name, v.Type, v.Md5, v.Created, v.Modified, uint64(v.Size), ri, nil)
	}
	ri.resources = resources
	return ri, nil
}

func (yd *yandexDisk) ReadResource(path string) ([]byte, error) {
	yd.client.GetResourceDownloadLink(path, nil)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, e
		}
		panic(err)
	}
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
