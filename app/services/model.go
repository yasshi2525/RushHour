package services

import (
	"fmt"
	"reflect"
	"strings"
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

	owner, _ := FetchOwner(msg.OName)

	switch msg.Op {
	case "create":
		rv := reflect.ValueOf(mkFuncs[msg.Target])
		switch msg.Target {
		case entities.PLAYER:
			level := entities.Normal
			if strings.Compare(msg.OName, "admin") == 0 {
				level = entities.Admin
			}
			CreatePlayer(msg.OName, msg.OName, msg.OName, level)
		case entities.RESIDENCE:
			fallthrough
		case entities.COMPANY:
			rv.Call([]reflect.Value{
				reflect.ValueOf(owner),
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
		case entities.RAILNODE:
			result := rv.Call([]reflect.Value{
				reflect.ValueOf(owner),
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
			if result[1].IsNil() {
				ExtendRailNode(owner, result[0].Interface().(*entities.RailNode), msg.X+10, msg.Y+10)
			}
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
			msg.ID, ok = randID(msg.Target, owner)
			if !ok {
				revel.AppLog.Warnf("no deleting data %s", msg.Target)
				break
			}
		}
		rv.Call([]reflect.Value{
			reflect.ValueOf(owner),
			reflect.ValueOf(msg.ID)})
	}
}

// randID return random id existing in repository
func randID(t entities.ModelType, owner *entities.Player) (uint, bool) {
	mapdata := Meta.Map[t]
	for _, key := range mapdata.MapKeys() {
		if mapdata.MapIndex(key).Interface().(entities.Ownable).Permits(owner) {
			return uint(key.Uint()), true
		}
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
