package services

import (
	"time"

	"github.com/revel/revel"
)

var gamemaster *time.Ticker

// StartProcedure start game.
func StartProcedure() {
	gamemaster = time.NewTicker(Config.Game.Interval.Duration)

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
	defer WarnLongExec(start, 2, "ゲーム進行", false)
	MuRoute.Lock()
	defer MuRoute.Unlock()
	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	time.Sleep(600 * time.Millisecond)
}
