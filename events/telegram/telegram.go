package telegram

import "awesomeProject4/cleints/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
}

func New(tg *telegram.Client, offset int) *Processor {}
