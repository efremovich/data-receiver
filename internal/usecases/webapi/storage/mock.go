package storage

import (
	"context"
	"fmt"
	"os"
	"path"
	"sync"
)

type storageClientMockImpl struct {
	dir string
	m   *sync.Mutex
}

func NewMockStorageClient(pathToDir string) (Storage, error) {
	_ = os.Mkdir(pathToDir, 0600)

	info, err := os.Stat(pathToDir)
	if err != nil {
		return nil, err
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("%s - не директория", pathToDir)
	}

	return &storageClientMockImpl{
		dir: pathToDir,
		m:   &sync.Mutex{},
	}, nil
}

func (s *storageClientMockImpl) SaveFile(_ context.Context, fileName string, data []byte) error {
	return os.WriteFile(path.Join(s.dir, fileName), data, 0600)
}

func (s *storageClientMockImpl) GetFile(_ context.Context, fileName string) ([]byte, error) {
	return os.ReadFile(path.Join(s.dir, fileName))
}

func (s storageClientMockImpl) Ping(_ context.Context) error {
	return nil
}
