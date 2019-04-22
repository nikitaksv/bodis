package storage

import "time"

type AuthData map[string]string

type Storage interface {
	Info() (Info, error)
	ReadResource(path string) (Resource, error)
	WriteResource(path string, data []byte) error
	DeleteResource(path string) error
	StateResource(path string)
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

type Resource interface {
	Path() string
	Name() string
	IsDir() bool
	// File or Directory size (bytes).
	Size() uint64

	Hash() string
	Permissions() Permissions
	Created() time.Duration
	Modified() time.Duration

	ParentResource() *Resource
}

type Permissions interface {
	IsWrite() bool
	IsRead() bool
	IsDelete() bool
}

type State struct {
	Exists bool
}
