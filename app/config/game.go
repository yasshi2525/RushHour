package config

import (
	"time"
)

type duration struct {
	D time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.D, err = time.ParseDuration(string(text))
	return err
}

// CnfResidence is configuration about residence
type CnfResidence struct {
	Interval  duration
	Capacity  int     `validate:"gt=0"`
	Randomize float64 `validate:"gte=0"`
}

type CnfCompany struct {
	Attract float64 `validate:"gt=0"`
}

// CnfGate is configuration about gate
type CnfGate struct {
	Num int `validate:"gt=0"`
}

// CnfPlatform is configuration about platform
type CnfPlatform struct {
	Capacity  int     `validate:"gt=0"`
	Randomize float64 `validate:"gte=0"`
}

// CnfTrain is configuration about train
type CnfTrain struct {
	Speed     float64 `validate:"gt=0"`
	Capacity  int     `validate:"gt=0"`
	Mobility  int     `validate:"gt=0"`
	Slowness  float64 `validate:"gt=0,lte=1"`
	Randomize float64 `validate:"gte=0"`
}

// CnfHuman is configuration about human
type CnfHuman struct {
	Speed float64 `validate:"gt=0"`
}

// CnfEntity is entity section of game.conf
type CnfEntity struct {
	MaxScale  float64 `toml:"max_scale" validate:"gtfield=MinScale"`
	MinScale  float64 `toml:"min_scale" validate:"ltfield=MaxScale"`
	Residence CnfResidence
	Company   CnfCompany
	Gate      CnfGate
	Platform  CnfPlatform
	Train     CnfTrain
	Human     CnfHuman
}

// CnfProcedure is configuration about game proceding
type CnfProcedure struct {
	Interval   duration
	Queue      uint `validate:"gt=0"`
	Simulation bool
}

// CnfRouting is configuration about paralization
type CnfRouting struct {
	Worker int `validate:"gt=0"`
	Alert  int `validate:"gte=0"`
}

// CnfBackup is configuration about backup interval
type CnfBackup struct {
	Interval duration
}

// CnfPerf is configuration about performance logging
type CnfPerf struct {
	View      duration
	Game      duration
	Operation duration
	Routing   duration
	Backup    duration
	Restore   duration
	Init      duration
}

// CnfService is service section of game.conf
type CnfService struct {
	Procedure CnfProcedure
	Routing   CnfRouting
	Backup    CnfBackup
	Perf      CnfPerf
}

// CnfGame is root element of game.conf
type CnfGame struct {
	Entity  CnfEntity
	Service CnfService
}
