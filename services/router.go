package services

import (
	"context"
	"log"
	"time"

	"github.com/yasshi2525/RushHour/entities"

	"github.com/yasshi2525/RushHour/route"
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
	alertEnabled := conf.Game.Service.Routing.Alert > 0

	lock, template, ok := scan(ctx)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= conf.Game.Service.Routing.Alert {
			log.Printf("routing was canceled (1/3) in scanning phase by %d times", routingBlockConunt)
		}
		return
	}

	payload, ok := search(ctx, template)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= conf.Game.Service.Routing.Alert {
			log.Printf("routing was canceled (2/3) in searching phase (%d/%d) by %d times",
				payload.Processed, payload.Total, routingBlockConunt)
		}
		return
	}

	RouteTemplate = payload
	reflectModel()
	if alertEnabled && routingBlockConunt >= conf.Game.Service.Routing.Alert {
		log.Printf("routing was successfully ended after %d times blocking", routingBlockConunt)
	}
	routingBlockConunt = 0
	if initialRouting { // force log when first routing after reboot
		WarnLongExec(start, lock, conf.Game.Service.Perf.Routing.D, "initial routing", true)
	} else {
		WarnLongExec(start, lock, conf.Game.Service.Perf.Routing.D, "routing")
	}
}

func scan(ctx context.Context) (time.Time, *route.Model, bool) {
	MuModel.RLock()
	defer MuModel.RUnlock()

	lock := time.Now()
	model, ok := route.Scan(ctx, Model)
	return lock, model, ok
}

func search(ctx context.Context, template *route.Model) (*route.Payload, bool) {
	return route.Search(ctx, entities.COMPANY, conf.Game.Service.Routing.Worker, template)
}

func reflectModel() {
	MuModel.Lock()
	defer MuModel.Unlock()

	for _, model := range RouteTemplate.Route {
		for _, n := range model.Nodes[entities.HUMAN] {
			Model.Humans[n.ID].Current = Model.Steps[n.ViaEdge.ID]
		}
	}
}
