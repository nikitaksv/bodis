package config

var configInstance *TomlConfig

type TomlConfig struct {
	Main    Main
	Storage struct {
		YandexDisk YandexDisk
	}
}

type Main struct {
	LogLevel       uint8
	ListenAddr     string
	ListenPort     uint
	DefaultTimeout uint64
}
type storage struct {
	Enabled bool
}
type YandexDisk struct {
	storage
	Token string
}

type FileError struct {
	Message  string
	FilePath string
	Table    string
	Key      string
	Value    string
	Line     int
}

func (e FileError) Error() string {
	return "[" + e.Table + "] Key=" + e.Key + " Value=" + e.Value + " File=" + e.FilePath + "#" + string(e.Line) + " Message=" + e.Message
}
