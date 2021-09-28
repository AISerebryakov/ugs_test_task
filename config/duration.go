package config

import (
	"strings"
	"time"
)

type Duration struct {
	time.Duration
}

// ParseDuration Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
func ParseDuration(durationStr string) (d Duration, err error) {
	if err := d.parse(durationStr); err != nil {
		return Duration{}, err
	}
	return d, nil
}

func (d Duration) TimeDuration() time.Duration {
	return d.Duration
}

func (d *Duration) UnmarshalYAML(f func(interface{}) error) error {
	var durationStr string
	err := f(&durationStr)
	if err != nil {
		return err
	}
	if err := d.parse(durationStr); err != nil {
		return err
	}
	return nil
}

func (d *Duration) parse(durationStr string) (err error) {
	if len(durationStr) == 0 {
		return nil
	}
	durationStr = strings.Replace(durationStr, " ", "", -1)
	durationStr = strings.ToLower(durationStr)
	d.Duration, err = time.ParseDuration(durationStr)
	if err != nil {
		return err
	}
	return nil
}
