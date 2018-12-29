package services

import (
	"time"

	"github.com/revel/revel"
)

var gamemaster *time.Ticker

// StartProcedure start game.
func StartProcedure() {
	gamemaster = time.NewTicker(1 * time.Second)

	go proceed()
}

func proceed() {
	for range gamemaster.C {
		start := time.Now()

		// 経路探索中の場合、ゲームを進行しない
		MuRoute.Lock()

		MuStatic.Lock()

		MuDynamic.Lock()

		time.Sleep(600 * time.Millisecond)

		MuDynamic.Unlock()
		MuStatic.Unlock()

		MuRoute.Unlock()

		WarnLongExec(start, 2, "ゲーム進行", false)
	}
}

// StopProcedure stop game
func StopProcedure() {
	if gamemaster != nil {
		revel.AppLog.Info("中止処理 開始")
		gamemaster.Stop()
		revel.AppLog.Info("中止処理 終了")
	}
}
