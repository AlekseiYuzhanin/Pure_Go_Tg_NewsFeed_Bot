package telegram

import (
	err2 "awesomeProject4/lib/err"
	"awesomeProject4/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatId int, username string) error {
	text = strings.TrimSpace(text)

	log.Printf("got new command: %s from %s", text, username)

	if isAddCmd(text) {

	}

	switch text {
	case RndCmd:
	case HelpCmd:
	case StartCmd:
	default:
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = err2.WrapIfErr(err, "cant do command: save page") }()
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return p.tg.SendMessage(chatID, "Already exists")
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}
