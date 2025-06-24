package app

import (
	"github.com/shivamhw/piro-go/pkg/log"
	tele "gopkg.in/telebot.v4"
)


func Run() error {
	b, err := NewBot()
	if err != nil {
		return err
	}
	b.Use(RecoverMiddleware)
	b.Handle(&searchBtn, b.searchBtnHandler)
	b.Handle(&scrapeBtn, b.scrpCallbackHandler)
	b.Handle(&subBtn, b.subBtnCallbackHandle)
	b.Handle("/scrape", b.scrapeCmdHandler) 
	b.Handle("/start", func(c tele.Context) error {
		return c.Send(welcomeMsg)
	})
	b.Handle("/search", b.searchCmdHandler)
	b.Handle(tele.OnText, b.defaultMsgHandler)
	log.Debugf("Bot started.....")
	b.Start()
	return nil
}
