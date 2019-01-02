package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// Operation represents request for model
type Operation struct {
	Source string
	Target entities.StaticRes
	ID     uint
	X      float64
	Y      float64
	Op     string
	OName  string
}

var modelChannel chan *Operation

var mkFuncs map[entities.StaticRes]interface{}
var rmFuncs map[entities.StaticRes]interface{}

// StartModelWatching setup watching model
func StartModelWatching() {

	modelChannel = make(chan *Operation, Config.Game.Queue)

	mkFuncs = make(map[entities.StaticRes]interface{})
	mkFuncs[entities.PLAYER] = CreatePlayer
	mkFuncs[entities.RESIDENCE] = CreateResidence
	mkFuncs[entities.COMPANY] = CreateCompany
	mkFuncs[entities.RAILNODE] = CreateRailNode

	rmFuncs = make(map[entities.StaticRes]interface{})
	rmFuncs[entities.RESIDENCE] = RemoveResidence
	rmFuncs[entities.COMPANY] = RemoveCompany
	rmFuncs[entities.RAILNODE] = RemoveRailNode

	go watchModel()
	revel.AppLog.Info("model watching was successfully started.")
}

// StopModelWatching closes channel
func StopModelWatching() {
	if modelChannel != nil {
		close(modelChannel)
		revel.AppLog.Info("model watching was successfully stopped.")
	}
}

func watchModel() {
	for msg := range modelChannel {
		start := time.Now()

		CancelRouting(msg.Source)

		processMsg(msg)
		StartRouting(msg.Source)

		WarnLongExec(start, 2, fmt.Sprintf("モデル変更(%v)", msg), false)
	}
	revel.AppLog.Info("model watching channel was closed.")
}

func processMsg(msg *Operation) {
	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	switch msg.Op {
	case "create":
		rv := reflect.ValueOf(mkFuncs[msg.Target])
		switch msg.Target {
		case entities.PLAYER:
			CreatePlayer(msg.OName, msg.OName, msg.OName, entities.Normal)
		case entities.RESIDENCE:
			fallthrough
		case entities.COMPANY:
			rv.Call([]reflect.Value{
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
		default:
			if owner, err := FetchOwner(msg.OName); err == nil {
				rv.Call([]reflect.Value{
					reflect.ValueOf(owner),
					reflect.ValueOf(msg.X),
					reflect.ValueOf(msg.Y)})
			} else {
				revel.AppLog.Warnf("invalid Player: %s", err)
			}
		}

	case "remove":
		rv := reflect.ValueOf(rmFuncs[msg.Target])
		if msg.ID == 0 {
			var ok bool
			msg.ID, ok = randID(msg.Target)
			if !ok {
				revel.AppLog.Warnf("no deleting data %s", msg.Target)
				break
			}
		}
		switch msg.Target {
		case entities.RESIDENCE:
			fallthrough
		case entities.COMPANY:
			rv.Call([]reflect.Value{reflect.ValueOf(msg.ID)})
		default:
			if owner, err := FetchOwner(msg.OName); err == nil {
				rv.Call([]reflect.Value{
					reflect.ValueOf(owner),
					reflect.ValueOf(msg.ID)})
			} else {
				revel.AppLog.Warnf("invalid Player: %s", err)
			}
		}
	}
}

// randID return random id existing in repository
func randID(t entities.StaticRes) (uint, bool) {
	mapdata := Meta.StaticMap[t]
	for _, e := range mapdata.MapKeys() {
		return uint(e.Uint()), true
	}
	revel.AppLog.Warnf("nodata %s", t)
	return 0, false
}

// UpdateModel queues user request.
func UpdateModel(msg *Operation) {
	revel.AppLog.Infof("updatemodel op = %+v", *msg)
	select {
	case modelChannel <- msg:
	default:
		revel.AppLog.Errorf("モデル変更キュー溢れ %+v", *msg)
	}
}
