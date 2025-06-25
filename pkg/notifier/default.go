package notifier

import "fmt"

type DefaultNtfy struct{}

func (d *DefaultNtfy) Notify(s string) {
	fmt.Printf("NTFY: %s\n", s)
}