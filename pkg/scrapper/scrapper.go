package scrapper

import (

	scr "github.com/shivamhw/content-pirate/pkg/scrapper"
	"github.com/shivamhw/content-pirate/store"
	"github.com/shivamhw/content-pirate/sources"
)

type Scrapper struct {
	s *scr.ScrapperV1
}

func StartScrapper() (*Scrapper, error) {
	cfg := scr.ScrapeCfg{
	TopicWorkers: 1,
	ImgWorkers: 1,
	}
	s, err := scr.NewScrapper(&cfg)
	if err != nil {
		return nil, err
	}
	go s.Start()
	return &Scrapper{
		s : s,
	}, nil
}

func (s *Scrapper) CheckJob(id string) (scr.TaskStatus, error) {
	return s.s.CheckJob(id)
}

func (s *Scrapper) GetJob(id string) (scr.Task, error) {
	return s.s.GetJob(id)
}

func (s *Scrapper) Scrape(src string, dst *store.DstPath, opts *sources.ScrapeOpts) (string, error) {
	j := scr.Job{
		SrcAc: src,
		Dst: *dst,
		Opts: *opts,
		SourceStore: scr.REDDIT,
	}
	return s.s.SubmitJob(j)
}
