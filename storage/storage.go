package storage

type Storage struct {
	filename string
}

func (s *Storage) GetFilename() string {
	return s.filename
}
