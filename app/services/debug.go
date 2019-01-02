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
	removeInterval = 1 * time.Minute
)

type opCallback func(source string, target entities.StaticRes)

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
	for _, target := range []entities.StaticRes{entities.RESIDENCE, entities.COMPANY} {
		tickOp("admin", target, updateInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(mkOp(src, tar))
		})
		tickOp("admin", target, removeInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(rmOp(src, tar))
		})
	}

	// user
	for i := 0; i < numUser; i++ {
		source := fmt.Sprintf("user%d", i)
		UpdateModel(mkOp(source, entities.PLAYER))

		tickOp(source, entities.RAILNODE, updateInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(mkOp(src, tar))
		})

		tickOp(source, entities.RAILNODE, removeInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(rmOp(src, tar))
		})
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
func mkOp(src string, target entities.StaticRes) *Operation {
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
func rmOp(src string, target entities.StaticRes) *Operation {
	return &Operation{
		Source: src,
		Op:     "remove",
		Target: target,
		OName:  src,
	}
}

func tickOp(source string, target entities.StaticRes, interval time.Duration, callback opCallback) {
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
func WarnLongExec(start time.Time, max float64, title string, verbose bool) {
	if consumed := time.Now().Sub(start).Seconds(); consumed > max {
		revel.AppLog.Warnf("%s に %.1f sec 消費", title, consumed)
	} else if verbose {
		revel.AppLog.Debugf("%s に %.1f sec 消費", title, consumed)
	}
}
