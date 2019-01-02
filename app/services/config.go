package services

import (
	"fmt"
	"time"

	"github.com/revel/revel"

	"github.com/BurntSushi/toml"
	validator "gopkg.in/go-playground/validator.v9"
)

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
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

type cfgTrain struct {
	Weight float64 `validate:"gt=0"`
}

type cfgHuman struct {
	Weight float64 `validate:"gt=0"`
}

type cfgGame struct {
	Interval duration
	Queue    uint `validate:"gt=0"`
}

type cfgBackup struct {
	Interval duration
}

type config struct {
	Residence cfgResidence
	Company   cfgCompany
	Train     cfgTrain
	Human     cfgHuman
	Game      cfgGame
	Backup    cfgBackup
}

// Config defines game feature
var Config config

// LoadConf load and validate game.conf
func LoadConf() {
	if _, err := toml.DecodeFile("conf/game.conf", &Config); err != nil {
		panic(fmt.Errorf("failed to load conf: %v", err))
	}

	if err := validator.New().Struct(Config); err != nil {
		panic(err)
	}
	revel.AppLog.Info("config file was successfully loaded.")
}
