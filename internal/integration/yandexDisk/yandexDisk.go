package yandexDisk

import (
	"bytes"
	"context"
	"github.com/nikitaksv/bodis/pkg/errors"
	"net/http"
	"strings"

	"github.com/nikitaksv/yandex-disk-sdk-go"

	"github.com/nikitaksv/bodis/pkg/storage"
)

var StorageKey = "yandexdisk"

type yandexDisk struct {
	client yadisk.YaDisk
}

func NewYandexDisk(ctx context.Context, client *http.Client, token string) *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(ctx, client, &yadisk.Token{AccessToken: token})
	if err != nil {
		panic(err) // fatal error
	}

	return &yandexDisk{client: yadClient}
}

func (yd *yandexDisk) Info() (storage.Info, errors.Error) {
	disk, err := yd.client.GetDisk([]string{"total_space", "max_file_size", "used_space"})
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@Info()": {},
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@Info()": {},
		})
	}
	return newYandexDiskInfo(uint64(disk.TotalSpace), uint64(disk.UsedSpace), uint64(disk.MaxFileSize)), nil
}

func (yd *yandexDisk) GetResourceInfo(path string) (storage.ResourceInfo, errors.Error) {
	return yd.getResourceInfo(convertPath(path))
}

func (yd *yandexDisk) getResourceInfo(path string) (*resourceInfo, errors.Error) {
	path = convertPath(path)
	// TODO Придумать как обойти limit и offset
	res, err := yd.client.GetResource(path, nil, 100, 0, false, "", "size")
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@GetResource()": {"path": path},
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@GetResource()": {"path": path},
		})
	}

	ri := newResourceInfo(yd, res.Path, res.Name, res.Type, res.Md5, res.Created, res.Modified, uint64(res.Size), nil, nil)
	resources := make([]resourceInfo, len(res.Embedded.Items))
	for i, v := range res.Embedded.Items {
		resources[i] = *newResourceInfo(yd, v.Path, v.Name, v.Type, v.Md5, v.Created, v.Modified, uint64(v.Size), ri, nil)
	}
	ri.resources = resources
	return ri, nil
}

func (yd *yandexDisk) ReadResource(path string) ([]byte, errors.Error) {
	path = convertPath(path)
	_, err := yd.client.GetResourceDownloadLink(path, nil)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@GetResourceDownloadLink()": {"path": path},
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@GetResourceDownloadLink()": {"path": path},
		})
	}
	// TODO придумать как скачивать файл. Возможно нужно реализовать в СДК
	panic("implement me")
}

func (yd *yandexDisk) WriteResource(path string, resInfo storage.ResourceInfo, data *bytes.Buffer) errors.Error {
	path = convertPath(path)
	uploadLink, err := yd.client.GetResourceUploadLink(path, nil, true)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@GetResourceUploadLink()": {"path": path, "data": map[string]interface{}{
					"len": data.Len(),
					"cap": data.Cap(),
				}},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@GetResourceUploadLink()": {"path": path, "data": map[string]interface{}{
				"len": data.Len(),
				"cap": data.Cap(),
			}},
		})
	}
	_, err = yd.client.PerformPartialUpload(uploadLink, data, calcPartSize(data.Len()))
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@PerformPartialUpload()": {"path": path, "data": map[string]interface{}{
					"len": data.Len(),
					"cap": data.Cap(),
				}},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@PerformPartialUpload()": {"path": path, "data": map[string]interface{}{
				"len": data.Len(),
				"cap": data.Cap(),
			}},
		})
	}
	return nil
}

func (yd *yandexDisk) DeleteResource(path string) errors.Error {
	path = convertPath(path)
	_, err := yd.client.DeleteResource(path, nil, true, "", false)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]map[string]interface{}{
				"@DeleteResource()": {"path": path},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]map[string]interface{}{
			"@DeleteResource()": {"path": path},
		})
	}
	return nil
}

func convertPath(path string) string {
	if !strings.HasPrefix(path, "disk:/") {
		if strings.HasPrefix(path, "/") {
			path = "disk:" + path
		} else {
			path = "disk:/" + path
		}
	}
	return path
}

func calcPartSize(dataLen int) int64 {
	count := dataLen / int(yadisk.MaxFileUploadSize)
	if count > 0 {
		return int64(dataLen / count)
	}
	return int64(dataLen)
}
