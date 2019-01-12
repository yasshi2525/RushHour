package services

import (
	"time"

	"github.com/revel/revel"
)

var gamemaster *time.Ticker

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
	start := time.Now()
	defer WarnLongExec(start, Const.Perf.Game.D, "procedure")
	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()
}
