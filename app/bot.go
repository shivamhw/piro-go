package app

import (
	"context"
	"os"
	"time"

	"github.com/shivamhw/content-pirate/pkg/kv"
	"github.com/shivamhw/content-pirate/pkg/reddit"
	"github.com/shivamhw/piro-go/pkg/log"
	"github.com/shivamhw/piro-go/pkg/notifier"
	"github.com/shivamhw/piro-go/pkg/scrapper"
	"gopkg.in/telebot.v4"
)

const (
	TOKEN_ENV_VAR = "TELE_TOKEN"
	NTFY_TOPIC = "NTFY_TOPIC"
)

type Bot struct {
	b   *telebot.Bot
	r   *reddit.RedditClient
	s   *scrapper.Scrapper
	ctx context.Context
	kv  kv.KV
	n   notifier.Notifier
}

func NewBot() (*Bot, error) {
	var n notifier.Notifier
	r, err := reddit.NewRedditClient(context.Background(), reddit.RedditClientOpts{})
	if err != nil {
		return nil, err
	}
	token := os.Getenv(TOKEN_ENV_VAR)
	topic := os.Getenv(NTFY_TOPIC)
	if token == "" {
		log.Errorf("env var not set", "var", TOKEN_ENV_VAR)
	}
	if topic == "" {
		log.WarnF("using default notifier")
		n = &notifier.DefaultNtfy{}
	} else {
		log.Infof("Using topic to ntfy", "topic", topic)
		n = notifier.NewNtfy(topic)
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
		n:   n,
		kv:  kv.GetInMemoryKv(),
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


func (b *Bot) SendAlbum(to telebot.Recipient, a telebot.Album, opts ...interface{}) (m []telebot.Message, err error) {
	m, err = b.b.SendAlbum(to, a, opts...)
	if err != nil {
		log.Errorf("sending album failed, reverting to single msg", "err", err)
		for _, f := range a {
			if f.MediaType() == "photo"{
				fiel := telebot.Photo{
					File: *f.MediaFile(),
				}
				msg, err :=  b.b.Send(to, &fiel)
				if err != nil  {
					log.Errorf(err.Error())
				} else {
					m = append(m, *msg)
				}
			}
		}
	}
	return m, nil
}

func (b *Bot) Notify(s string) {
	b.n.Notify(s)
}