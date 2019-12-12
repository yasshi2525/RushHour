package services

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/yasshi2525/RushHour/entities"
)

var (
	numUser        = 5
	viewInterval   = 10 * time.Second
	updateInterval = 120 * time.Second
	removeInterval = 120 * time.Second
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
	UpdateModel(mkOp("admin", entities.PLAYER))

	// admin
	for _, target := range []entities.ModelType{entities.RESIDENCE, entities.COMPANY} {
		tickOp("admin", target, updateInterval, func(src string, tar entities.ModelType) {
			UpdateModel(mkOp(src, tar))
		})
		tickOp("admin", target, removeInterval, func(src string, tar entities.ModelType) {
			UpdateModel(rmOp(src, tar))
		})
		UpdateModel(&Operation{
			Source: "admin",
			Op:     "create",
			Target: entities.RESIDENCE,
			OName:  "admin",
			X:      0,
			Y:      0,
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
			entities.LINETASK,
			entities.TRAIN} {
			tickOp(source, target, updateInterval, func(src string, tar entities.ModelType) {
				UpdateModel(mkOp(src, tar))
			})
			tickOp(source, target, removeInterval, func(src string, tar entities.ModelType) {
				UpdateModel(rmOp(src, tar))
			})
		}
	}

	op := mkOp("user0", entities.RAILNODE)
	op.X, op.Y = 0, 0
	UpdateModel(op)

	simWg.Wait()
	log.Println("simulation was successfully started.")
}

// StopSimulation stop simulation.
func StopSimulation() {
	if simCh != nil {
		close(simCh)
		for _, t := range simTickers {
			t.Stop()
		}
		log.Println("simulation was successfully stopped.")
	}
}

func watchSim() {
	for v := range simCh {
		simTickers = append(simTickers, v)
	}
	log.Println("simulation channel was successfully stopped.")
}

// mkOp returns creation operation
func mkOp(src string, target entities.ModelType) *Operation {
	size := 1 << (conf.Game.Entity.MaxScale - conf.Game.Entity.MinScale)
	return &Operation{
		Source: src,
		Op:     "create",
		Target: target,
		OName:  src,
		X:      float64(size) * rand.Float64(),
		Y:      float64(size) * rand.Float64(),
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
func WarnLongExec(start time.Time, lock time.Time, max time.Duration, title string, verbose ...bool) {
	if consumed := time.Now().Sub(start); consumed > max {
		log.Printf("%s consumed %.2f(%.2f) sec >%.2f", title, consumed.Seconds(),
			time.Now().Sub(lock).Seconds(), max.Seconds())
	} else if len(verbose) > 0 && verbose[0] {
		log.Printf("%s consumed %.2f(%.2f) sec <%.2f", title, consumed.Seconds(),
			time.Now().Sub(lock).Seconds(), max.Seconds())
	}
}
