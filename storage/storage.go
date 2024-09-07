package storage

import (
	err2 "awesomeProject4/lib/err"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

var ErrNoSavedPages = errors.New("No saved pages")

type Page struct {
	URL      string
	UserName string
}

func (p *Page) Hash() (string, error) {
	h := sha1.New()

	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", err2.Wrap(err, "cant calculate hash")
	}

	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", err2.Wrap(err, "cant calculate hash")
	}

	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
