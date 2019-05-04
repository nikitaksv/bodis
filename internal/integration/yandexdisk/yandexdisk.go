package yandexdisk

import (
	"bytes"
	"context"
	"net/http"
	"strings"

	"github.com/nikitaksv/bodis/pkg/errors"

	"github.com/nikitaksv/yandex-disk-sdk-go"

	"github.com/nikitaksv/bodis/pkg/storage"
)

var StorageKey storage.IntegrationKey = "yandexdisk"

type yandexDisk struct {
	client yadisk.YaDisk
}

func New(ctx context.Context, client *http.Client, token string) *yandexDisk {
	yadClient, err := yadisk.NewYaDisk(ctx, client, &yadisk.Token{AccessToken: token})
	if err != nil {
		panic(err) // fatal error
	}

	return &yandexDisk{client: yadClient}
}

func (yd *yandexDisk) Info() (storage.Info, storage.Error) {
	disk, err := yd.client.GetDisk([]string{"total_space", "max_file_size", "used_space"})
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]interface{}{
				"method": "Info()",
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method": "Info()",
		})
	}
	return newYandexDiskInfo(uint64(disk.TotalSpace), uint64(disk.UsedSpace), uint64(disk.MaxFileSize)), nil
}

func (yd *yandexDisk) GetResourceInfo(path string) (storage.ResourceInfo, storage.Error) {
	return yd.getResourceInfo(convertPath(path))
}

func (yd *yandexDisk) getResourceInfo(path string) (*resourceInfo, storage.Error) {
	path = convertPath(path)
	// TODO Придумать как обойти limit и offset
	res, err := yd.client.GetResource(path, nil, 100, 0, false, "", "size")
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]interface{}{
				"method":    "GetResource()",
				"arguments": map[string]string{"path": path},
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method":    "GetResource()",
			"arguments": map[string]string{"path": path},
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

func (yd *yandexDisk) ReadResource(path string) (*bytes.Buffer, storage.Error) {
	path = convertPath(path)
	_, err := yd.client.GetResourceDownloadLink(path, nil)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return nil, newError(e.ErrorID, e.Description, map[string]interface{}{
				"method":    "GetResourceDownloadLink()",
				"arguments": map[string]string{"path": path},
			})
		}
		return nil, newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method":    "GetResourceDownloadLink()",
			"arguments": map[string]string{"path": path},
		})
	}
	// TODO придумать как скачивать файл. Возможно нужно реализовать в СДК
	panic("implement me")
}

func (yd *yandexDisk) WriteResource(path string, resInfo storage.ResourceInfo, data *bytes.Buffer) storage.Error {
	path = convertPath(path)
	uploadLink, err := yd.client.GetResourceUploadLink(path, nil, true)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]interface{}{
				"method": "GetResourceUploadLink()",
				"arguments": map[string]interface{}{
					"path": path,
					"data": map[string]int{
						"len": data.Len(),
						"cap": data.Cap(),
					},
				},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method": "GetResourceUploadLink()",
			"arguments": map[string]interface{}{
				"path": path,
				"data": map[string]int{
					"len": data.Len(),
					"cap": data.Cap(),
				},
			},
		})
	}
	_, err = yd.client.PerformPartialUpload(uploadLink, data, calcPartSize(data.Len()))
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]interface{}{
				"method": "PerformPartialUpload()",
				"arguments": map[string]interface{}{
					"path": path,
					"data": map[string]int{
						"len": data.Len(),
						"cap": data.Cap(),
					}},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method": "PerformPartialUpload()",
			"arguments": map[string]interface{}{
				"path": path,
				"data": map[string]int{
					"len": data.Len(),
					"cap": data.Cap(),
				}},
		})
	}
	return nil
}

func (yd *yandexDisk) DeleteResource(path string) storage.Error {
	path = convertPath(path)
	_, err := yd.client.DeleteResource(path, nil, true, "", false)
	if err != nil {
		if e, ok := err.(*yadisk.Error); ok {
			return newError(e.ErrorID, e.Description, map[string]interface{}{
				"method":    "DeleteResource()",
				"arguments": map[string]string{"path": path},
			})
		}
		return newError(errors.InternalSDK, err.Error(), map[string]interface{}{
			"method":    "DeleteResource()",
			"arguments": map[string]string{"path": path},
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
