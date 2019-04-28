package storage

type Error interface {
	error
	StorageID() string
	ErrorID() string
	Description() string
	Params() map[string]string
}

type BaseError struct {
	StorageID string
}

func (e *BaseError) Error() string {
	return "[" + e.StorageID() + "] " + "(" + e.ErrorID() + ")" + e.Description()
}

func (e *BaseError) StorageID() string {
	panic("implement me")
}

func (e *BaseError) ErrorID() string {
	panic("implement me")
}

func (e *BaseError) Description() string {
	panic("implement me")
}

func (e *BaseError) Params() map[string]string {
	panic("implement me")
}
