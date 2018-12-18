package services

import (
	"sync"
	"time"

	"github.com/revel/revel"
)

var gamemaster *time.Ticker
var modelLock *sync.Mutex

func Start() {
	modelLock = new(sync.Mutex)
	gamemaster = time.NewTicker(1 * time.Second)

	go func() {
		for range gamemaster.C {
			modelLock.Lock()
			revel.AppLog.Info("ゲーム 開始")
			time.Sleep(5 * time.Second)
			revel.AppLog.Info("ゲーム 終了")
			modelLock.Unlock()
		}
	}()

	go func() {
		time.Sleep(55 * time.Second)
		Stop()
	}()

	go func() {
		for i := 0; i < 10; i++ {
			modelLock.Lock()
			revel.AppLog.Info("経路探索 開始")
			time.Sleep(5 * time.Second)
			revel.AppLog.Info("経路探索 終了")
			modelLock.Unlock()
			time.Sleep(1 * time.Second)
		}
	}()
}

func Stop() {
	if gamemaster != nil {
		revel.AppLog.Info("中止処理 開始")
		gamemaster.Stop()
		revel.AppLog.Info("中止処理 終了")
	}
}
