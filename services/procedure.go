package services

import (
	"log"
	"time"
)

var gamemaster *time.Ticker
var beforeProcedure time.Time

// StartProcedure start game.
func StartProcedure() {
	gamemaster = time.NewTicker(conf.Game.Service.Procedure.Interval.D)

	go watchGame()
	log.Println("game procedure was successfully started.")
}

// StopProcedure stop game
func StopProcedure() {
	if gamemaster != nil {
		gamemaster.Stop()
		log.Println("game procedure was successfully stopped.")
	}
}

func watchGame() {
	for range gamemaster.C {
		processGame()
	}
}

func processGame() {
	start := time.Now()
	MuModel.Lock()
	defer MuModel.Unlock()
	lock := time.Now()
	defer WarnLongExec(start, lock, conf.Game.Service.Perf.Game.D, "procedure")
	defer func() { beforeProcedure = time.Now() }()

	if beforeProcedure.IsZero() {
		return
	}
	interval := time.Now().Sub(beforeProcedure).Seconds()

	for _, t := range Model.Trains {
		t.Step(interval)
	}
}
