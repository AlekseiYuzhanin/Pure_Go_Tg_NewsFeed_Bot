package files

import (
	err2 "awesomeProject4/lib/err"
	"awesomeProject4/storage"
	"encoding/gob"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

const (
	defaultPerm = 0774
)

var ErrNoSavedPages = errors.New("No saved pages")

type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath}
}

func (storage *Storage) Save(p *storage.Page) (err error) {
	defer func() { err = err2.WrapIfErr(err, "cant save page") }()

	filePath := filepath.Join(storage.basePath, p.UserName)

	if err = os.MkdirAll(filePath, defaultPerm); err != nil {
		return err
	}

	fName, err := fileName(p)
	if err != nil {
		return err
	}
	filePath = filepath.Join(filePath, fName)

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gob.NewEncoder(file).Encode(p); err != nil {
		return err
	}

	return nil
}

func (s *Storage) PickRandom(userName string) (page *storage.Page, err error) {
	defer func() { err = err2.WrapIfErr(err, "cant pick random page") }()

	fPath := filepath.Join(s.basePath, userName)

	files, err := os.ReadDir(fPath)
	if err != nil {
		return nil, err
	}

	if len(files) == 0 {
		return nil, ErrNoSavedPages
	}

	rand.Seed(time.Now().UTC().UnixNano())
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(fPath, file.Name()))
}

func (s *Storage) Remove(page *storage.Page) (err error) {
	fileName, err := fileName(page)
	if err != nil {
		return err2.Wrap(err, "cant remove page")
	}
	path := filepath.Join(s.basePath, page.UserName, fileName)
	if err := os.Remove(path); err != nil {
		return err2.Wrap(err, fmt.Sprintf("can't remove file %s", path))
	}
	return nil
}

func (s *Storage) IsExists(page *storage.Page) (bool, error) {
	filename, err := fileName(page)
	if err != nil {
		return false, err2.Wrap(err, "cant remove page")
	}
	path := filepath.Join(s.basePath, page.UserName, filename)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, err2.Wrap(err, "can't check if file exists")
	}

	return true, nil
}

func (s *Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err2.Wrap(err, "Can't open file")
	}
	defer f.Close()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, err2.Wrap(err, "Can't decode file")
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
