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
	MuRoute.Lock()
	defer MuRoute.Unlock()

	start := time.Now()
	defer WarnLongExec(start, Config.Perf.Routing.D, "routing")

	alertEnabled := Config.Routing.Alert > 0

	template, ok := scan(ctx)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= Config.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (1/3) in scanning phase by %d times", routingBlockConunt)
		}
		return
	}

	payload, ok := search(ctx, template)
	if !ok {
		routingBlockConunt++
		if alertEnabled && routingBlockConunt >= Config.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (2/3) in searching phase (%d/%d) by %d times",
				payload.Processed, payload.Total, routingBlockConunt)
		}
		return
	}

	RouteTemplate = payload
	reflectModel()
	if alertEnabled && routingBlockConunt > Config.Routing.Alert {
		revel.AppLog.Infof("routing was successfully ended after %d times blocking", routingBlockConunt)
	}
	routingBlockConunt = 0
}

func scan(ctx context.Context) (*route.Model, bool) {
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	return route.Scan(ctx, Model, Meta)
}

func search(ctx context.Context, template *route.Model) (*route.Payload, bool) {
	return route.Search(ctx, entities.COMPANY, Config.Routing.Worker, template)
}

func reflectModel() {
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	for _, a := range Model.Agents {
		for _, n := range RouteTemplate.Route[a.Human.ToID].Nodes {
			if n.SameAs(a.Human) {
				a.Current = Model.Steps[n.ViaEdge.ID]
			}
		}
	}
}
