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
	Capacity  uint    `validate:"min=1"`
	Randomize float64 `validate:"min=0"`
}

type cfgCompany struct {
	Scale float64 `validate:"gt=0"`
}

type cfgGate struct {
	Num uint `validate:"gt=0"`
}

type cfgPlatform struct {
	Capacity uint `validate:"gt=0"`
}

type cfgTrain struct {
	Weight   float64 `validate:"gt=0"`
	Slowness float64 `validate:"gt=0,lte=1"`
}

type cfgHuman struct {
	Weight float64 `validate:"gt=0"`
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
