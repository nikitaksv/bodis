package yandexDisk

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/nikitaksv/bodis/pkg/storage"
)

type diskInfo struct {
	totalSpace      uint64
	usedSpace       uint64
	freeSpace       uint64
	maxResourceSize uint64
	isRemote        bool
}

func newYandexDiskInfo(totalSpace uint64, usedSpace uint64, maxResourceSize uint64) *diskInfo {
	return &diskInfo{
		totalSpace:      totalSpace,
		usedSpace:       usedSpace,
		freeSpace:       totalSpace - usedSpace,
		maxResourceSize: maxResourceSize,
		isRemote:        true,
	}
}

func (ydi *diskInfo) TotalSpace() uint64 {
	return ydi.totalSpace
}

func (ydi *diskInfo) UsedSpace() uint64 {
	return ydi.usedSpace
}

func (ydi *diskInfo) FreeSpace() uint64 {
	return ydi.freeSpace
}

func (ydi *diskInfo) MaxResourceSize() uint64 {
	return ydi.maxResourceSize
}

func (ydi *diskInfo) IsRemote() bool {
	return ydi.isRemote
}

type resourceInfo struct {
	client *yandexDisk

	ext         string
	path        string
	name        string
	isDir       bool
	size        uint64
	hash        string
	permissions permissions
	created     time.Time
	modified    time.Time

	parentResource *resourceInfo
	resources      []resourceInfo
}

func newResourceInfo(client *yandexDisk, path, name, dir, hash, createdS, modifiedS string, size uint64, parentResource *resourceInfo, resources []resourceInfo) *resourceInfo {
	ext := filepath.Ext(name)
	isDir := false
	if dir == "dir" {
		isDir = true
		ext = ""
	}

	created, err := time.Parse(time.RFC3339, createdS)
	if err != nil {
		created = time.Time{}
	}
	modified, err := time.Parse(time.RFC3339, modifiedS)
	if err != nil {
		modified = time.Time{}
	}

	if err != nil {
		fmt.Println(err)
	}

	return &resourceInfo{
		client:         client,
		ext:            ext,
		path:           path,
		name:           name,
		isDir:          isDir,
		size:           size,
		hash:           hash,
		permissions:    permissions{},
		created:        created,
		modified:       modified,
		parentResource: parentResource,
		resources:      resources,
	}
}

func (ri resourceInfo) Extension() string {
	return ri.ext
}

func (ri resourceInfo) Path() string {
	return ri.path
}

func (ri resourceInfo) Name() string {
	return ri.name
}

func (ri resourceInfo) IsDir() bool {
	return ri.isDir
}

func (ri resourceInfo) Size() uint64 {
	return ri.size
}

func (ri resourceInfo) Hash() string {
	return ri.hash
}

func (ri resourceInfo) Permissions() storage.Permissions {
	return &ri.permissions
}

func (ri resourceInfo) Created() time.Time {
	return ri.created
}

func (ri resourceInfo) Modified() time.Time {
	return ri.modified
}

func (ri resourceInfo) ParentResource() storage.ResourceInfo {
	return ri.parentResource
}

func (ri resourceInfo) Resources() []storage.ResourceInfo {
	models := make([]storage.ResourceInfo, len(ri.resources))
	for i, v := range ri.resources {
		models[i] = resourceInfo(v)
	}

	if ri.isDir && len(models) == 0 {
		newRi, err := ri.client.getResourceInfo(ri.path)
		if err != nil {
			return nil
		}
		if len(newRi.resources) == 0 {
			return nil
		}
		newRi.parentResource = &ri
		return newRi.Resources()
	}

	return models
}

// Owner has all rights.
type permissions struct {
}

func (permissions) IsWrite() bool {
	return true
}

func (permissions) IsRead() bool {
	return true
}
