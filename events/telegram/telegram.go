package telegram

import (
	"awesomeProject4/cleints/telegram"
	"awesomeProject4/events"
	err2 "awesomeProject4/lib/err"
	"awesomeProject4/storage"
	"errors"
)

var ErrUnknown = errors.New("unknown event")
var ErrUnknownMeta = errors.New("unknown meta")

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(tg *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      tg,
		offset:  0,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, err2.Wrap(err, "cant get events")
	}

	if len(update) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(update))

	for _, update := range update {
		res = append(res, event(update))
	}

	p.offset = update[len(update)-1].ID + 1

	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return ErrUnknown
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return err2.Wrap(err, "cant process message")
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return err2.Wrap(err, "cant process message")
	}
	return nil
}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, ErrUnknownMeta
	}
	return res, nil
}

func event(update telegram.Update) events.Event {
	updateType := fetchType(update)
	res := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   update.Message.Chat.ID,
			Username: update.Message.From.Username,
		}
	}

	return res
}

func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}
	return update.Message.Text
}

func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}

	return events.Message
}
