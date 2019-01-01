package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var (
	numUser        = 10
	viewInterval   = 1 * time.Second
	updateInterval = 10 * time.Second
	backupInterval = 30 * time.Second
)

type opCallback func(source string, target entities.StaticRes)

// Main immitates some user requests some actions.
// TODO remove it
func Main() {
	Restore()
	StartModelWatching()
	StartProcedure()

	// admin
	for _, target := range []entities.StaticRes{entities.RESIDENCE, entities.COMPANY} {
		tickOp("admin", target, updateInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(mkOp(src, tar))
		})
		tickOp("admin", target, updateInterval, func(src string, tar entities.StaticRes) {
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

		tickOp(source, entities.RAILNODE, updateInterval, func(src string, tar entities.StaticRes) {
			UpdateModel(rmOp(src, tar))
		})
	}

	var backup = time.NewTicker(backupInterval)
	go func() {
		for range backup.C {
			Backup()
		}
	}()
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
	go func() {
		sleep := rand.Intn(int(interval.Seconds()))
		time.Sleep(time.Duration(sleep) * time.Second)
		t := time.NewTicker(interval)
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
