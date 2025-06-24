package app

import "gopkg.in/telebot.v4"

const (
	SCRP_BTN   = "scrp_btn"
	SUB_BTN    = "sub_btn"
	SEARCH_BTN = "search_btn"
)

var (
	selector  = &telebot.ReplyMarkup{}
	menu      = &telebot.ReplyMarkup{ResizeKeyboard: true}
	scrapeBtn = selector.Data("", SCRP_BTN, "")
	subBtn    = selector.Data("", SUB_BTN, "")
	searchBtn = selector.Data("", SEARCH_BTN, "")
	aboutBtn  = menu.Text("ℹ️ About me")
)
