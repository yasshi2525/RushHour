package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var (
	numUser        = 5
	viewInterval   = 1 * time.Second
	updateInterval = 30 * time.Second
	removeInterval = 120 * time.Second
)

type opCallback func(source string, target entities.ModelType)

var simCh chan *time.Ticker
var simTickers []*time.Ticker
var simWg *sync.WaitGroup

func AfterStart() {
	p, _ := CreatePlayer("test", "test", "test", entities.Normal)
	rn1, _ := CreateRailNode(p, 10, 10)
	st, _ := CreateStation(p, rn1, "test1")
	l, _ := CreateRailLine(p, "test1", true)
	StartRailLine(p, l, st.Platform)
	rn2, e1, _, _ := ExtendRailNode(p, rn1, 20, 20)
	InsertLineTaskRailEdge(p, e1, false)
	_, e2, _, _ := ExtendRailNode(p, rn2, 30, 30)
	InsertLineTaskRailEdge(p, e2, false)
	Backup()
}

// StartSimulation immitates some user requests some actions.
// TODO remove it
func StartSimulation() {
	simTickers = []*time.Ticker{}
	simCh = make(chan *time.Ticker)
	simWg = &sync.WaitGroup{}
	go watchSim()
	UpdateModel(mkOp("admin", entities.PLAYER))

	// admin
	for _, target := range []entities.ModelType{entities.RESIDENCE, entities.COMPANY} {
		tickOp("admin", target, updateInterval, func(src string, tar entities.ModelType) {
			UpdateModel(mkOp(src, tar))
		})
		tickOp("admin", target, removeInterval, func(src string, tar entities.ModelType) {
			UpdateModel(rmOp(src, tar))
		})
	}

	// user
	for i := 0; i < numUser; i++ {
		source := fmt.Sprintf("user%d", i)
		UpdateModel(mkOp(source, entities.PLAYER))

		for _, target := range []entities.ModelType{
			entities.RAILNODE,
			entities.RAILEDGE,
			entities.STATION,
			entities.RAILLINE,
			entities.LINETASK} {
			tickOp(source, target, updateInterval, func(src string, tar entities.ModelType) {
				UpdateModel(mkOp(src, tar))
			})
			tickOp(source, target, removeInterval, func(src string, tar entities.ModelType) {
				UpdateModel(rmOp(src, tar))
			})
		}

	}
	simWg.Wait()
	revel.AppLog.Info("simulation was succeesfully started.")
}

// StopSimulation stop simulation.
func StopSimulation() {
	if simCh != nil {
		close(simCh)
		for _, t := range simTickers {
			t.Stop()
		}
		revel.AppLog.Info("simulation was succeesfully stopped.")
	}
}

func watchSim() {
	for v := range simCh {
		simTickers = append(simTickers, v)
	}
	revel.AppLog.Info("simulation channel was succeesfully stopped.")
}

// mkOp returns creation operation
func mkOp(src string, target entities.ModelType) *Operation {
	return &Operation{
		Source: src,
		Op:     "create",
		Target: target,
		OName:  src,
		X:      rand.Float64() * 100,
		Y:      rand.Float64() * 100,
	}
}

// rmOp returns deletion operation
func rmOp(src string, target entities.ModelType) *Operation {
	return &Operation{
		Source: src,
		Op:     "remove",
		Target: target,
		OName:  src,
	}
}

func tickOp(source string, target entities.ModelType, interval time.Duration, callback opCallback) {
	simWg.Add(1)
	go func() {
		simWg.Done()
		sleep := time.Duration(rand.Intn(int(interval.Seconds())))
		time.Sleep(sleep * time.Second)
		t := time.NewTicker(interval)
		simCh <- t
		for range t.C {
			callback(source, target)
		}
	}()
}

// WarnLongExec alerts long time consuming task.
func WarnLongExec(start time.Time, max time.Duration, title string, verbose ...bool) {
	if consumed := time.Now().Sub(start); consumed > max {
		revel.AppLog.Warnf("%s consumed %.2f sec >%.2f", title, consumed.Seconds(), max.Seconds())
	} else if len(verbose) > 0 && verbose[0] {
		revel.AppLog.Debugf("%s consumed %.2f sec <%.2f", title, consumed.Seconds(), max.Seconds())
	}
}