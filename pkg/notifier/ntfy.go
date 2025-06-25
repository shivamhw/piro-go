package notifier

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shivamhw/piro-go/pkg/log"
)

const (
	BaseUrl = "http://ntfy.sh/"
)

type Ntfy struct {
	topic string
}

func NewNtfy(t string) *Ntfy {
	return &Ntfy{
		topic: t,
	}
}

func (n *Ntfy) Notify(s string){
	url := fmt.Sprintf("%s%s", BaseUrl, n.topic)
	_, err := http.Post(url, "", strings.NewReader(s))
	if err != nil {
		log.Errorf("in ntfy ", "err", err)
	}
}