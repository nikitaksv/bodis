package storage

import (
	"bytes"
	"time"
)

type IntegrationKey string
type AuthData map[string]string

type Storage interface {
	Info() (Info, Error)
	GetResourceInfo(path string) (ResourceInfo, Error)
	//TODO возможно придется поменять сигнатуру.
	ReadResource(path string) (*bytes.Buffer, Error)
	WriteResource(path string, resInfo ResourceInfo, data *bytes.Buffer) Error
	DeleteResource(path string) Error
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

// Error interface for storage.
// This error is used for logging.
type Error interface {
	error
	// Storage key (exp. "yandexdisk", "googledrive")
	StorageKey() IntegrationKey
	// Error ID (exp. "401 Unauthorized")
	ErrorID() string
	// Description (exp. "Invalid token")
	Description() string
	// Params
	//
	// Map key is the namespace (method name)
	//
	// "Info()":{ "token": "123456", "data": { "length": "123456" }, ... }
	Params() map[string]interface{}
}

// BaseError implement Error interface.
type BaseError struct {
	storageKey  IntegrationKey
	errorID     string
	description string
	params      map[string]interface{}
}

func NewBaseError(storageKey IntegrationKey, errorID string, description string, params map[string]interface{}) *BaseError {
	return &BaseError{
		storageKey:  storageKey,
		errorID:     errorID,
		description: description,
		params:      params,
	}
}
func (e BaseError) Error() string {
	return "[" + string(e.StorageKey()) + "] " + "(" + e.ErrorID() + ") \"" + e.Description() + "\""
}

func (e BaseError) StorageKey() IntegrationKey {
	return e.storageKey
}

func (e BaseError) ErrorID() string {
	return e.errorID
}

func (e BaseError) Description() string {
	return e.description
}

func (e BaseError) Params() map[string]interface{} {
	return e.params
}
