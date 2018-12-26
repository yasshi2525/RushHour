package services

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/revel/revel"
)

func Main() {
	Restore()
	StartModelWatching()
	StartProcedure()

	for i := 0; i < 10; i++ {
		go func() {
			userView := time.NewTicker(1 * time.Second)
			for range userView.C {
				ViewMap()
			}
		}()

		go func(id int) {
			time.Sleep(time.Duration(rand.Intn(30)) * time.Second)
			userOp := time.NewTicker(30 * time.Second)
			revel.AppLog.Debugf("start operation %d", id)
			for range userOp.C {
				ViewMap()
				ChangeMap(fmt.Sprintf("user%d", id))
			}
		}(i)
	}

	/*
		var backup = time.NewTicker(1 * time.Minute)
		go func() {
			for range backup.C {
				Backup()
			}
		}()*/
}

func WarnLongExec(start time.Time, max float64, title string, verbose bool) {
	if consumed := time.Now().Sub(start).Seconds(); consumed > max {
		revel.AppLog.Warnf("%s に %.1f sec 消費", title, consumed)
	} else if verbose {
		revel.AppLog.Debugf("%s に %.1f sec 消費", title, consumed)
	}
}
