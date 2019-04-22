package integration

import (
	"context"
	"github.com/nikitaksv/bodis/pkg/storage"
	"github.com/nikitaksv/yandex-disk-sdk-go"
	"net/http"
)

type yandexDisk struct {
	client yadisk.YaDisk
}

func newYandexDisk(ctx context.Context, client http.Client, token string) *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(ctx, &client, &yadisk.Token{AccessToken: token})
	if err != nil {
		panic(err.Error())
	}
	return &yandexDisk{client: yadClient}
}

func (yd *yandexDisk) Info() (storage.Info, error) {
	disk, err := yd.client.GetDisk(nil)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			// TODO Log error
			return nil, e
		}
		panic(err.Error())
	}
	return newYandexDiskInfo(uint64(disk.TotalSpace), uint64(disk.UsedSpace), uint64(disk.MaxFileSize)), nil
}

func (yd *yandexDisk) ReadResource(path string) (storage.Resource, error) {
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

type yandexDiskInfo struct {
	totalSpace      uint64
	usedSpace       uint64
	freeSpace       uint64
	maxResourceSize uint64
	isRemote        bool
}

func newYandexDiskInfo(totalSpace uint64, usedSpace uint64, maxResourceSize uint64) *yandexDiskInfo {
	return &yandexDiskInfo{
		totalSpace:      totalSpace,
		usedSpace:       usedSpace,
		freeSpace:       totalSpace - usedSpace,
		maxResourceSize: maxResourceSize,
		isRemote:        true,
	}
}

func (ydi *yandexDiskInfo) TotalSpace() uint64 {
	return ydi.totalSpace
}

func (ydi *yandexDiskInfo) UsedSpace() uint64 {
	return ydi.usedSpace
}

func (ydi *yandexDiskInfo) FreeSpace() uint64 {
	return ydi.freeSpace
}

func (ydi *yandexDiskInfo) MaxResourceSize() uint64 {
	return ydi.maxResourceSize
}

func (ydi *yandexDiskInfo) IsRemote() bool {
	return ydi.isRemote
}
