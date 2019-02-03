package services

import (
	"time"

	"github.com/revel/revel"
)

var gamemaster *time.Ticker
var beforeProcedure time.Time

// StartProcedure start game.
func StartProcedure() {
	gamemaster = time.NewTicker(Const.Game.Interval.D)

	go watchGame()
	revel.AppLog.Info("game procedure was successfully started.")
}

// StopProcedure stop game
func StopProcedure() {
	if gamemaster != nil {
		gamemaster.Stop()
		revel.AppLog.Info("game procedure was successfully stopped.")
	}
}

func watchGame() {
	for range gamemaster.C {
		processGame()
	}
}

func processGame() {
	revel.AppLog.Debug("[DEBUG] process Model after UnLock")
	start := time.Now()
	revel.AppLog.Debug("[DEBUG] process Model before Lock")
	MuModel.Lock()
	revel.AppLog.Debug("[DEBUG] process Model after Lock")
	defer MuModel.Unlock()
	lock := time.Now()
	defer WarnLongExec(start, lock, Const.Perf.Game.D, "procedure")
	defer func() { beforeProcedure = time.Now() }()

	if beforeProcedure.IsZero() {
		return
	}
	interval := time.Now().Sub(beforeProcedure).Seconds()

	for _, t := range Model.Trains {
		t.Step(interval)
	}
}
