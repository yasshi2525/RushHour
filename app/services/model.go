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
	Target entities.ModelType
	ID     uint
	X      float64
	Y      float64
	Op     string
	OName  string
}

var modelChannel chan *Operation

var mkFuncs map[entities.ModelType]interface{}
var rmFuncs map[entities.ModelType]interface{}

// StartModelWatching setup watching model
func StartModelWatching() {

	modelChannel = make(chan *Operation, Config.Game.Queue)

	mkFuncs = make(map[entities.ModelType]interface{})
	mkFuncs[entities.PLAYER] = CreatePlayer
	mkFuncs[entities.RESIDENCE] = CreateResidence
	mkFuncs[entities.COMPANY] = CreateCompany
	mkFuncs[entities.RAILNODE] = CreateRailNode
	mkFuncs[entities.STATION] = CreateStation

	rmFuncs = make(map[entities.ModelType]interface{})
	rmFuncs[entities.RESIDENCE] = RemoveResidence
	rmFuncs[entities.COMPANY] = RemoveCompany
	rmFuncs[entities.RAILNODE] = RemoveRailNode
	rmFuncs[entities.STATION] = RemoveStation

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
		skipReroute := msg.Target == entities.PLAYER

		if !skipReroute {
			CancelRouting(msg.Source)
		}
		processMsg(msg)
		if !skipReroute {
			StartRouting(msg.Source)
		}
		WarnLongExec(start, Config.Perf.Operation.D, fmt.Sprintf("operation(%v)", msg))
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
		owner, _ := FetchOwner(msg.OName)
		switch msg.Target {
		case entities.PLAYER:
			CreatePlayer(msg.OName, msg.OName, msg.OName, entities.Normal)
		case entities.RESIDENCE:
			fallthrough
		case entities.COMPANY:
			rv.Call([]reflect.Value{
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
		case entities.RAILNODE:
			rv.Call([]reflect.Value{
				reflect.ValueOf(owner),
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
		case entities.STATION:
			if rn := randRailNode(owner); rn != nil {
				rv.Call([]reflect.Value{
					reflect.ValueOf(owner),
					reflect.ValueOf(rn),
					reflect.ValueOf("NoName")})
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
func randID(t entities.ModelType) (uint, bool) {
	mapdata := Meta.Map[t]
	for _, e := range mapdata.MapKeys() {
		return uint(e.Uint()), true
	}
	revel.AppLog.Warnf("nodata %s", t)
	return 0, false
}

// UpdateModel queues user request.
func UpdateModel(msg *Operation) {
	//revel.AppLog.Infof("updatemodel op = %+v", *msg)
	select {
	case modelChannel <- msg:
	default:
		revel.AppLog.Errorf("モデル変更キュー溢れ %+v", *msg)
	}
}

func randRailNode(o *entities.Player) *entities.RailNode {
	for _, rn := range Model.RailNodes {
		if rn.Permits(o) {
			return rn
		}
	}
	return nil
}
