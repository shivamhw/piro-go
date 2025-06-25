package app

import (
	"fmt"
	"time"

	"github.com/shivamhw/content-pirate/pkg/scrapper"
	"github.com/shivamhw/content-pirate/sources"
	"github.com/shivamhw/piro-go/pkg/log"
	tele "gopkg.in/telebot.v4"
)

func (b *Bot) scrape(update *tele.Message, red string, opts *sources.ScrapeOpts) ([]sources.Item, error) {
	deadline := time.Now().Add(10 * time.Minute) 
	var files []sources.Item
	id, err := b.s.Scrape(red, dst, opts)
	if err != nil {
		return nil, err
	}
	for  {
		st, err := b.s.CheckJob(id)
		if err != nil {
			return files, err
		}
		b.b.Edit(update, fmt.Sprintf("Downloading %d/%d", st.ItemDone, st.TotalItem))
		if st.ItemDone >= st.TotalItem && st.Status != scrapper.TaskCreated{
			res, err := b.s.GetJob(id)
			if err != nil {
				return files, err
			}
			files = res.I
			break
		}
		if time.Now().After(deadline) {
			log.Errorf("timeout while executing the task, not waiting anymore %s", id)
			return files, nil
		}
		time.Sleep(2 * time.Second)
	}
	return files, nil
}

func (b *Bot) sendScrapped(c tele.Context, files []sources.Item) error {
	al := tele.Album{}
	for _, file := range files {
		path := fmt.Sprintf("%s/%s/%s", dst.BasePath, dst.ImgPath, file.Dst)
		log.Debugf("sending %s to %s", path, c.Sender().FirstName)
		f := &tele.Photo{
			File:    tele.FromDisk(path),
			Caption: file.Title,
		}
		al = append(al, f)
		if len(al) > 9 {
			_, err := b.SendAlbum(c.Sender(), al)
			if err != nil {
				log.Errorf("error sending msg", err.Error())
			}
			al = al[:0]
		}
	}
	if len(al) > 0 {
		_, err := b.SendAlbum(c.Sender(), al)
		if err != nil {
			return err
		}
	}
	return nil
}

func (b *Bot) search(ctx tele.Context, q string) error {
	log.Debugf("searching for %s", q)
	ctx.Send(fmt.Sprintf("searching for %s", q))
	res, err := b.r.SearchSubreddits(q, 10)
	if err != nil {
		return err
	}
	inline := &tele.ReplyMarkup{}

	var rows []tele.Row
	for _, r := range res {
		btn := inline.Data(r.Name, SUB_BTN, r.Name)
		rows = append(rows, inline.Row(btn))
	}
	inline.Inline(rows...)
	ctx.Send(fmt.Sprintf("Found %d results:", len(res)), inline)
	return nil
}
