package services

import (
	"context"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/models"
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
		models.MuRoute.Lock()
		defer models.MuRoute.Unlock()
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
	models.MuStatic.RLock()
	defer models.MuStatic.RUnlock()

	models.MuAgent.RLock()
	defer models.MuAgent.RUnlock()

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
	models.MuAgent.Lock()
	defer models.MuAgent.Unlock()

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
