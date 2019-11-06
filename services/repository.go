package services

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/route"
)

// Model is contained all data for gaming
var Model *entities.Model

// RouteTemplate is default route information in order to avoid huge calculation.
var RouteTemplate *route.Payload

// MuModel is mutex lock for Model
var MuModel sync.RWMutex

// MuRoute is mutex lock for routing
var MuRoute sync.Mutex

type OpLog struct {
	gorm.Model
	Op        string
	OwnerID   uint
	Obj1      string
	Obj2      string
	Obj3      string
	Obj4      string
	idx       uint
	TimeStamp time.Time
}

func (op *OpLog) Add(obj entities.Entity) {
	str := fmt.Sprintf("%s(%d)", obj.B().Type().Short(), obj.B().Idx())
	switch op.idx {
	case 0:
		op.Obj1 = str
	case 1:
		op.Obj2 = str
	case 2:
		op.Obj3 = str
	case 3:
		op.Obj4 = str
	default:
		log.Printf("too many args = %d", op.idx+1)
	}
	op.idx++
}

func AddOpLog(op string, o *entities.Player, args ...entities.Entity) {
	log := &OpLog{Op: op, OwnerID: o.ID, TimeStamp: time.Now()}
	for _, obj := range args {
		log.Add(obj)
	}
	OpCache = append(OpCache, log)
}

var OpCache []*OpLog

// InitLock must prepare first.
func InitLock() {
	MuModel = sync.RWMutex{}
	MuRoute = sync.Mutex{}
}

// InitRepository initialize storage
func InitRepository() {
	Model = entities.NewModel(serviceConf.AppConf.Game.Entity, serviceConf.Auther)
	OpCache = []*OpLog{}
}
