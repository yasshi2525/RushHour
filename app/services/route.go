package services

import (
	"context"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

var cancelChannel chan string
var routingContext context.Context
var routingCancel context.CancelFunc
var IsSearching bool

func StartRouting(msg string) {
	routingContext, routingCancel = context.WithCancel(context.Background())
	searchCtx, searchCancel := context.WithCancel(routingContext)
	reflectCtx, reflectCancel := context.WithCancel(routingContext)

	go func() {
		entities.MuRoute.Lock()
		defer entities.MuRoute.Unlock()
		defer routingCancel()
		search(searchCtx, searchCancel, msg)
		reflect(reflectCtx, reflectCancel, msg)
	}()
}

func CancelRouting(msg string) {
	if routingCancel != nil {
		routingCancel()
		revel.AppLog.Infof("Canceleld by %s", msg)
	} else {
		revel.AppLog.Warn("routingCancel is nil")
	}
}

func search(ctx context.Context, cancel context.CancelFunc, msg string) {
	defer cancel()
	entities.MuStatic.RLock()
	defer entities.MuStatic.RUnlock()

	entities.MuAgent.RLock()
	defer entities.MuAgent.RUnlock()

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			revel.AppLog.Infof("search Canceleld %s in %d / 10", msg, i+1)
			return
		default:
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

func reflect(ctx context.Context, cancel context.CancelFunc, msg string) {
	defer cancel()
	entities.MuAgent.Lock()
	defer entities.MuAgent.Unlock()

	for i := 0; i < 10; i++ {
		select {
		case <-ctx.Done():
			revel.AppLog.Infof("reflect Canceleld %s in %d / 10", msg, i+1)
			return
		default:
			time.Sleep(100 * time.Millisecond)
		}
	}
}
