package services

import (
	"fmt"
	"time"

	"github.com/revel/revel"
)

// Operation represents request for model
type Operation struct {
	Source string
	Target string
	ID     uint
	X      float64
	Y      float64
	Op     string
}

var modelChannel chan *Operation

// StartModelWatching setup watching model
func StartModelWatching() {
	defer revel.AppLog.Debug("モデルチャネル セットアップ終了")

	modelChannel = make(chan *Operation, 10)

	go watchModel()
}

func watchModel() {
	for msg := range modelChannel {
		start := time.Now()

		CancelRouting(msg.Source)

		MuStatic.Lock()
		MuDynamic.Lock()

		switch msg.Target {
		case "residence":
			switch msg.Op {
			case "create":
				CreateResidence(msg.X, msg.Y)
			case "remove":
				if msg.ID == 0 {
					for _, r := range Static.Residences {
						RemoveResidence(r.ID)
						break
					}
				} else {
					RemoveResidence(msg.ID)
				}
			}
		case "company":
			switch msg.Op {
			case "create":
				CreateCompany(msg.X, msg.Y)
			case "remove":
				if msg.ID == 0 {
					for _, c := range Static.Companies {
						RemoveCompany(c.ID)
						break
					}
				} else {
					RemoveCompany(msg.ID)
				}
			}
		}

		MuDynamic.Unlock()
		MuStatic.Unlock()

		StartRouting(msg.Source)

		WarnLongExec(start, 2, fmt.Sprintf("モデル変更(%v)", msg), false)
	}
}

// UpdateModel queues user request.
func UpdateModel(msg *Operation) {
	select {
	case modelChannel <- msg:
	default:
		revel.AppLog.Errorf("モデル変更キュー溢れ %v", msg)
	}
}
