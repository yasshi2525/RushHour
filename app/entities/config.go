package entities

import "time"

type duration struct {
	D time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.D, err = time.ParseDuration(string(text))
	return err
}

type cfgResidence struct {
	Interval  duration
	Capacity  int     `validate:"gt=0"`
	Randomize float64 `validate:"gte=0"`
}

type cfgCompany struct {
	Scale float64 `validate:"gt=0"`
}

type cfgGate struct {
	Num int `validate:"gt=0"`
}

type cfgPlatform struct {
	Capacity  int     `validate:"gt=0"`
	Randomize float64 `validate:"gte=0"`
}

type cfgTrain struct {
	Speed     float64 `validate:"gt=0"`
	Capacity  int     `validate:"gt=0"`
	Mobility  int     `validate:"gt=0"`
	Slowness  float64 `validate:"gt=0,lte=1"`
	Randomize float64 `validate:"gte=0"`
}

type cfgHuman struct {
	Speed  float64 `validate:"gt=0"`
}

type Config struct {
	Residence cfgResidence
	Company   cfgCompany
	Gate      cfgGate
	Platform  cfgPlatform
	Train     cfgTrain
	Human     cfgHuman
}

var Const Config
