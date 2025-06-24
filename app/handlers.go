package app

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/shivamhw/content-pirate/sources"
	"github.com/shivamhw/content-pirate/store"
	"github.com/shivamhw/piro-go/pkg/log"
	tele "gopkg.in/telebot.v4"
)

// todo : this should be placed somewhere else
var dst = &store.DstPath{
	BasePath: "./download",
	ImgPath:  "imgs",
	VidPath:  "vids",
}

func (b *Bot) scrpCallbackHandler(c tele.Context) error {
	red := strings.Split(c.Data(), "|")[0]
	lim := strings.Split(c.Data(), "|")[1]
	log.Infof("scrapping sub reddit with limits,", "sub", red, "limit", lim)
	limit, _ := strconv.Atoi(lim)
	opts := &sources.ScrapeOpts{
		Limit:      limit,
		Duration:   "month",
		SkipVideos: true,
	}
	msg, _ := b.Send(c.Sender(), fmt.Sprintf("started scrapping %s for %s posts", red, lim))
	files, err := b.scrape(red, opts)
	if err != nil {
		return err
	}
	b.b.Edit(msg, fmt.Sprintf("scrapped %d posts, sending it to you", len(files)))
	if err := b.sendScrapped(c, files); err != nil {
		return err
	}
	return c.Respond()
}

func (b *Bot) subBtnCallbackHandle(c tele.Context) error {
	mu := &tele.ReplyMarkup{}
	var btns []tele.Btn
	btns = append(btns, mu.Data("top 50", SCRP_BTN, c.Data(), "50"))
	btns = append(btns, mu.Data("top 25", SCRP_BTN, c.Data(), "25"))
	btns = append(btns, mu.Data("top 10", SCRP_BTN, c.Data(), "10"))
	row := mu.Row(btns...)
	mu.Inline(row)
	c.Send(fmt.Sprintf("Here are your options for %s", c.Data()), mu)
	return c.Respond()
}

func (b *Bot) scrapeCmdHandler(c tele.Context) error {
	red := strings.Split(c.Text(), " ")[1]
	opts := &sources.ScrapeOpts{
		Limit:      10,
		Duration:   "month",
		SkipVideos: true,
	}
	files, err := b.scrape(red, opts)
	if err != nil {
		return err
	}
	c.Send(fmt.Sprintf("scrapped %d posts, sending it to you", len(files)))
	if err := b.sendScrapped(c, files); err != nil {
		c.Send("somethings bad happedn")
		return err
	}
	return nil
}

func (b *Bot) searchCmdHandler(ctx tele.Context) error {
	q := strings.Join(strings.Split(ctx.Text(), " ")[1:], " ")
	return b.search(ctx, q)
}

func (b *Bot) defaultMsgHandler(ctx tele.Context) error {
	q := ctx.Text()
	log.Infof("got string returning back the search btn", "msg", q, "user", ctx.Sender().FirstName)
	in := &tele.ReplyMarkup{}
	btn := in.Data(fmt.Sprintf("🔍 Want to search for %s on reddit? 🔍", q), SEARCH_BTN, q)
	in.Inline(in.Row(btn))
	return ctx.Send("what do you want to do ????", in)
}

func (b *Bot) searchBtnHandler(ctx tele.Context) error {
	q := ctx.Data()
	log.Infof("searching for user", "q", q, "user", ctx.Sender().FirstName)
	ctx.Respond()
	return b.search(ctx, q)
}
