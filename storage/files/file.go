package files

import (
	err2 "awesomeProject4/lib/err"
	"awesomeProject4/storage"
	"os"
	"path/filepath"
)

const (
	defaultPerm = 0774
)

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath}
}

func (storage *Storage) Save(p *storage.Page) (err error) {
	defer func() { err = err2.WrapIfErr(err, "cant save") }()

	filePath := filepath.Join(storage.basePath, p.UserName)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

}
