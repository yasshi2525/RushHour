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
	updateInterval = 1 * time.Minute
	removeInterval = 1 * time.Minute
)

type opCallback func(source string, target entities.ModelType)

var simCh chan *time.Ticker
var simTickers []*time.Ticker
var simWg *sync.WaitGroup

// StartSimulation immitates some user requests some actions.
// TODO remove it
func StartSimulation() {
	simTickers = []*time.Ticker{}
	simCh = make(chan *time.Ticker)
	simWg = &sync.WaitGroup{}
	go watchSim()

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

		for _, target := range []entities.ModelType{entities.RAILNODE, entities.STATION} {
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
		revel.AppLog.Warnf("%s consumed %.2f sec (border %.2f sec)", title, consumed.Seconds(), max.Seconds())
	} else if len(verbose) > 0 && verbose[0] {
		revel.AppLog.Debugf("%s consumed %.2f sec", title, consumed.Seconds())
	}
}
