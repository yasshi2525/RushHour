package services

import (
	"fmt"

	"github.com/BurntSushi/toml"
	validator "gopkg.in/go-playground/validator.v9"
)

type residence struct {
	Interval  float64 `validate:"gt=0"`
	Capacity  uint    `validate:"min=1"`
	Randomize float64 `validate:"min=0"`
}

type company struct {
	Scale float64 `validate:"gt=0"`
}

type train struct {
	Weight float64 `validate:"gt=0"`
}

type human struct {
	Weight float64 `validate:"gt=0"`
}

type config struct {
	Residence residence
	Company   company
	Train     train
	Human     human
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
}
