package services

import (
	"context"
	"time"

	"github.com/revel/revel"
)

var cancelChannel chan string
var routingContext context.Context
var routingCancel context.CancelFunc

// IsSearching represents whether searching is executed or not.
var IsSearching bool

// StartRouting start searching.
func StartRouting(msg string) {
	routingContext, routingCancel = context.WithCancel(context.Background())
	searchCtx, searchCancel := context.WithCancel(routingContext)
	reflectCtx, reflectCancel := context.WithCancel(routingContext)

	go func() {
		MuRoute.Lock()
		defer MuRoute.Unlock()
		defer routingCancel()
		search(searchCtx, searchCancel, msg)
		reflectTo(reflectCtx, reflectCancel, msg)
	}()
}

// CancelRouting stop current executing searching.
func CancelRouting(msg string) {
	if routingCancel != nil {
		routingCancel()
		//revel.AppLog.Infof("Canceleld by %s", msg)
	} else {
		revel.AppLog.Warn("routingCancel is nil")
	}
}

func search(ctx context.Context, cancel context.CancelFunc, msg string) {
	defer cancel()
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			//revel.AppLog.Infof("search Canceleld %s in %d / 10", msg, i+1)
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func reflectTo(ctx context.Context, cancel context.CancelFunc, msg string) {
	defer cancel()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			//revel.AppLog.Infof("reflect Canceleld %s in %d / 10", msg, i+1)
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
