package services

import (
	"fmt"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/models"
)

var modelChannel chan string

func StartModelWatching() {
	defer revel.AppLog.Debug("モデルチャネル セットアップ終了")

	modelChannel = make(chan string, 10)

	go watchModel()
}

func watchModel() {
	for msg := range modelChannel {
		start := time.Now()

		CancelRouting(msg)

		models.MuStatic.Lock()
		models.MuAgent.Lock()

		models.MuAgent.Unlock()
		models.MuStatic.Unlock()

		StartRouting(msg)

		WarnLongExec(start, 2, fmt.Sprintf("モデル変更(%s)", msg), false)
	}
}

func UpdateModel(msg string) {
	select {
	case modelChannel <- msg:
	default:
		revel.AppLog.Errorf("モデル変更キュー溢れ %s", msg)
	}
}
