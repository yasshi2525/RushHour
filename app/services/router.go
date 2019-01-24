package services

import (
	"context"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services/route"
)

var routingContext context.Context
var routingCancel context.CancelFunc
var routingBlockConunt int

var searching bool

// IsSearching represents whether searching is executed or not.
func IsSearching() bool {
	return searching
}

func StartRouting() {
	if routingCancel != nil {
		routingCancel()
	}
	routingContext, routingCancel = context.WithCancel(context.Background())
	go processRouting(routingContext)
}

// CancelRouting stop current executing searching.
func CancelRouting() {
	if routingCancel != nil {
		routingCancel()
	}
}

func processRouting(ctx context.Context) {
	start := time.Now()
	MuRoute.Lock()
	defer MuRoute.Unlock()

	initialRouting := RouteTemplate == nil
	alertEnabled := Const.Routing.Alert > 0

	lock, template, ok := scan(ctx)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= Const.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (1/3) in scanning phase by %d times", routingBlockConunt)
		}
		return
	}

	payload, ok := search(ctx, template)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= Const.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (2/3) in searching phase (%d/%d) by %d times",
				payload.Processed, payload.Total, routingBlockConunt)
		}
		return
	}

	RouteTemplate = payload
	reflectModel()
	if alertEnabled && routingBlockConunt >= Const.Routing.Alert {
		revel.AppLog.Infof("routing was successfully ended after %d times blocking", routingBlockConunt)
	}
	routingBlockConunt = 0
	if initialRouting { // force log when first routing after reboot
		WarnLongExec(start, lock, Const.Perf.Routing.D, "initial routing", true)
	} else {
		WarnLongExec(start, lock, Const.Perf.Routing.D, "routing")
	}
}

func scan(ctx context.Context) (time.Time, *route.Model, bool) {
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	lock := time.Now()
	model, ok := route.Scan(ctx, Model)
	return lock, model, ok
}

func search(ctx context.Context, template *route.Model) (*route.Payload, bool) {
	return route.Search(ctx, entities.COMPANY, Const.Routing.Worker, template)
}

func reflectModel() {
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	for _, model := range RouteTemplate.Route {
		for _, n := range model.Nodes[entities.HUMAN] {
			Model.Humans[n.ID].Current = Model.Steps[n.ViaEdge.ID]
		}
	}
}
