package services

import (
	"fmt"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"

	"github.com/BurntSushi/toml"
	validator "gopkg.in/go-playground/validator.v9"
)

type duration struct {
	D time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.D, err = time.ParseDuration(string(text))
	return err
}

type cfgGame struct {
	Interval duration
	Queue    uint `validate:"gt=0"`
}

type cfgRouting struct {
	Worker int `validate:"gt=0"`
	Alert  int `validate:"gte=0"`
}

type cfgBackup struct {
	Interval duration
}

type cfgPerf struct {
	View      duration
	Game      duration
	Operation duration
	Routing   duration
	Backup    duration
	Restore   duration
	Init      duration
}

type cfgService struct {
	Game    cfgGame
	Routing cfgRouting
	Backup  cfgBackup
	Perf    cfgPerf
}

type config struct {
	Entity  entities.Config
	Service cfgService
}

// Config defines game feature
var Config config

// Const has constants for service.
var Const cfgService

// LoadConf load and validate game.conf
func LoadConf() {
	if _, err := toml.DecodeFile("conf/game.conf", &Config); err != nil {
		panic(fmt.Errorf("failed to load conf: %v", err))
	}
	if err := validator.New().Struct(Config); err != nil {
		panic(fmt.Errorf("%+v, %v", Config, err))
	}
	Const = Config.Service

	defer revel.AppLog.Info("config file was successfully loaded.")
	return
}
