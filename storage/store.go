package storage

type Store interface {
	Save(data []byte) error
	Load() ([]byte, error)
	GetFilename() string
}
