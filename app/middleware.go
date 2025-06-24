package app

import (
	"github.com/shivamhw/piro-go/pkg/log"
	tele "gopkg.in/telebot.v4"
)



func RecoverMiddleware(next tele.HandlerFunc) tele.HandlerFunc {
	return func(c tele.Context) error {
		defer func() {
			if r := recover(); r != nil {
				log.Infof("Recovered in handler: %v", r)
				_ = c.Send("Something went wrong, but I'm still alive!")
			}
		}()
		return next(c)
	}
}
