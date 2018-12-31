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
	OName  string
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
		case "player":
			CreatePlayer(msg.OName, msg.OName, msg.OName)
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
		case "rail_node":
			switch msg.Op {
			case "create":
				if o, err := FetchOwner(msg.OName); err == nil {
					CreateRailNode(o, msg.X, msg.Y)
				} else {
					revel.AppLog.Warnf("invalid Player: %s", err)
				}
			case "remove":
				if msg.ID == 0 {
					for _, rn := range Static.RailNodes {
						if o, err := FetchOwner(msg.OName); err == nil {
							RemoveRailNode(o, rn.ID)
						} else {
							revel.AppLog.Warnf("invalid Player: %s", err)
						}
						break
					}
				} else {
					if o, err := FetchOwner(msg.OName); err == nil {
						RemoveRailNode(o, msg.ID)
					} else {
						revel.AppLog.Warnf("invalid Player: %s", err)
					}
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
