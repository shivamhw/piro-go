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

func (s *Scrapper) Scrape(src string, dst *store.DstPath, opts *sources.ScrapeOpts) ([]sources.Item, error) {
	j := scr.Job{
		SrcAc: src,
		Dst: *dst,
		Opts: *opts,
		SourceStore: scr.REDDIT,
	}
	id, err := s.s.SubmitJob(j)
	if err != nil {
		return nil, err
	}
	s.s.WaitOnId(id)
	res, err := s.s.GetJob(id)
	if err !=nil {
		return nil, err
	}

	return res.I, nil
}
