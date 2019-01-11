package services

import (
	"fmt"
	"math/rand"
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
	mkFuncs[entities.RAILEDGE] = ExtendRailNode
	mkFuncs[entities.STATION] = CreateStation
	mkFuncs[entities.RAILLINE] = CreateRailLine

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
			CancelRouting()
		}
		processMsg(msg)
		if !skipReroute {
			StartRouting()
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
			rv.Call([]reflect.Value{
				reflect.ValueOf(owner),
				reflect.ValueOf(msg.X),
				reflect.ValueOf(msg.Y)})
		case entities.RAILEDGE:
			if raw := randEntity(owner, entities.RAILNODE); raw != nil {
				rv.Call([]reflect.Value{
					reflect.ValueOf(owner),
					reflect.ValueOf(raw),
					reflect.ValueOf(msg.X),
					reflect.ValueOf(msg.Y)})
			}
		case entities.STATION:
			if rn := randEntity(owner, entities.RAILNODE); rn != nil {
				rv.Call([]reflect.Value{
					reflect.ValueOf(owner),
					reflect.ValueOf(rn),
					reflect.ValueOf("NoName")})
			}
		case entities.RAILLINE:
			rv.Call([]reflect.Value{
				reflect.ValueOf(owner),
				reflect.ValueOf("NoName"),
				reflect.ValueOf(rand.Intn(2) == 0)})
		case entities.LINETASK:
			l := randEntity(owner, entities.RAILLINE)
			if l != nil {
				if p := randEntity(owner, entities.PLATFORM); p != nil {
					l, p := l.(*entities.RailLine), p.(*entities.Platform)
					StartRailLine(owner, l, p)
				}
				if re := randEntity(owner, entities.RAILEDGE); re != nil {
					re := re.(*entities.RailEdge)
					StartRailLineEdge(owner, l.(*entities.RailLine), re)
					InsertLineTaskRailEdge(owner, re, rand.Intn(2) == 0)
				}
				if rand.Intn(2) == 0 {
					CompleteRailLine(owner, l.(*entities.RailLine))
				}
			}
		}

	case "remove":
		rv := reflect.ValueOf(rmFuncs[msg.Target])
		if !rv.IsValid() {
			break
		}
		if msg.ID == 0 {
			var ok bool
			msg.ID, ok = randID(msg.Target, owner)
			if !ok {
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
	mapdata := Model.Values[t]
	for _, key := range mapdata.MapKeys() {
		if mapdata.MapIndex(key).Interface().(entities.Ownable).Permits(owner) {
			return uint(key.Uint()), true
		}
	}
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

func randEntity(o *entities.Player, res entities.ModelType) interface{} {
	var entity interface{}
	Model.ForEach(res, func(obj entities.Indexable) {
		if obj.(entities.Ownable).Permits(o) {
			entity = obj
		}
	})
	return entity
}
