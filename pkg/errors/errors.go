package errors

// Error interface for storage.
// This error is used for logging.
type Error interface {
	error
	// Storage key (exp. "yandexdisk", "googledrive")
	StorageKey() string
	// Error ID (exp. "401 Unauthorized")
	ErrorID() string
	// Description (exp. "Invalid token")
	Description() string
	// Params
	//
	// Map key is the namespace (method name)
	//
	// "Info()":{ "token": "123456", "data": { "length": "123456" }, ... }
	Params() map[string]map[string]interface{}
}

// BaseError implement Error interface.
type BaseError struct {
	storageKey  string
	errorID     string
	description string
	params      map[string]map[string]interface{}
}

func NewBaseError(storageKey string, errorID string, description string, params map[string]map[string]interface{}) *BaseError {
	return &BaseError{
		storageKey:  storageKey,
		errorID:     errorID,
		description: description,
		params:      params,
	}
}
func (e BaseError) Error() string {
	return "[" + e.StorageKey() + "] " + "(" + e.ErrorID() + ") \"" + e.Description() + "\""
}

func (e BaseError) StorageKey() string {
	return e.storageKey
}

func (e BaseError) ErrorID() string {
	return e.errorID
}

func (e BaseError) Description() string {
	return e.description
}

func (e BaseError) Params() map[string]map[string]interface{} {
	return e.params
}

var InternalSDK = "internal sdk"
