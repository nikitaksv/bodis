package storage

import "time"

type AuthData map[string]string

type Storage interface {
	Info() (Info, error)
	GetResourceInfo(path string) (ResourceInfo, error)
	ReadResource(path string) ([]byte, error)
	WriteResource(path string, resInfo ResourceInfo, data []byte) error
	DeleteResource(path string) error
}

type Info interface {
	// Total disk space (bytes).
	TotalSpace() uint64
	// Used disk space (bytes).
	UsedSpace() uint64
	// Free disk space (bytes).
	FreeSpace() uint64
	// Max write resource size (bytes).
	MaxResourceSize() uint64

	// Set to TRUE if it is a cloud storage (OneDrive, Yandex.Disk, Google Drive).
	IsRemote() bool
}

type ResourceInfo interface {
	Extension() string
	Path() string
	Name() string
	IsDir() bool
	// File or Directory size (bytes).
	Size() uint64

	Hash() string
	Permissions() Permissions
	Created() time.Time
	Modified() time.Time

	ParentResource() ResourceInfo
	Resources() []ResourceInfo
}

type Permissions interface {
	IsWrite() bool
	IsRead() bool
}

type State struct {
	Exists bool
}
