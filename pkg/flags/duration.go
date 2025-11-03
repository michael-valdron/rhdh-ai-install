package flags

import (
	"fmt"
	"time"

	"github.com/spf13/pflag"
)

type Duration struct {
	pflag.Value
	obj *time.Duration
	str string
}

func NewDuration(d *time.Duration) *Duration {
	return &Duration{obj: d}
}

func (d *Duration) Set(value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return fmt.Errorf("invalid duration '%s'", value)
	}
	*d.obj = duration
	d.str = value
	return nil
}

func (d *Duration) String() string {
	return d.str
}

func (d *Duration) Type() string {
	return "time.Duration"
}
