package app

import (
	"context"
	"os"
	"time"

	"github.com/shivamhw/content-pirate/pkg/reddit"
	"github.com/shivamhw/piro-go/pkg/log"
	"github.com/shivamhw/piro-go/pkg/scrapper"
	"gopkg.in/telebot.v4"
)

const (
	TOKEN_ENV_VAR = "TELE_TOKEN"
)

type Bot struct {
	b   *telebot.Bot
	r   *reddit.RedditClient
	s   *scrapper.Scrapper
	ctx context.Context
}

func NewBot() (*Bot, error) {
	r, err := reddit.NewRedditClient(context.Background(), reddit.RedditClientOpts{})
	if err != nil {
		return nil, err
	}
	token := os.Getenv(TOKEN_ENV_VAR)
	if token == "" {
		log.Errorf("env var not set", "var", TOKEN_ENV_VAR)
	}
	s, err := scrapper.StartScrapper()
	if err != nil {
		return nil, err
	}
	pref := telebot.Settings{
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
		Token:  token,
	}
	b, err := telebot.NewBot(pref)
	if err != nil {
		log.Errorf(err.Error())
		return nil, err
	}
	return &Bot{
		b:   b,
		r:   r,
		s:   s,
		ctx: context.Background(),
	}, nil
}

func (b *Bot) Use(t ...telebot.MiddlewareFunc) {
	b.b.Use(t...)
}

func (b *Bot) Handle(endpoint interface{}, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
	b.b.Handle(endpoint, h, m...)
}

func (b *Bot) Send(to telebot.Recipient, what interface{}, opts ...interface{}) (*telebot.Message, error) {
	return b.b.Send(to, what, opts...)
}

func (b *Bot) Start() {
	b.b.Start()
}


func (b *Bot) SendAlbum(to telebot.Recipient, a telebot.Album, opts ...interface{}) ([]telebot.Message, error) {
	return b.b.SendAlbum(to, a, opts...)
}