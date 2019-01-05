package services

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services/route"
)

var routingContext context.Context
var routingCancel context.CancelFunc
var routingBlockConunt int64

var searching bool

// IsSearching represents whether searching is executed or not.
func IsSearching() bool {
	return searching
}

func StartRouting() {
	routingContext, routingCancel = context.WithCancel(context.Background())
	go processRouting(routingContext)
}

// CancelRouting stop current executing searching.
func CancelRouting() {
	routingCancel()
	atomic.AddInt64(&routingBlockConunt, 1)
}

func processRouting(ctx context.Context) {
	start := time.Now()
	defer WarnLongExec(start, Config.Perf.Routing.D, "routing")

	template, ok := scan(ctx)
	if !ok {
		if b := int(routingBlockConunt); b >= Config.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (1/3) in scanning phase by %d times", b)
		}
		return
	}

	payload, ok := search(ctx, template)
	if !ok {
		if b := int(routingBlockConunt); b >= Config.Routing.Alert {
			revel.AppLog.Warnf("routing was canceled (2/3) in searching phase (%d/%d) by %d times",
				payload.Processed, payload.Total, b)
		}
		return
	}

	RouteTemplate = payload
	reflectModel()
	if b := int(routingBlockConunt); b >= Config.Routing.Alert {
		revel.AppLog.Infof("routing was successfully ended after %d times blocking", b)
	}
	atomic.StoreInt64(&routingBlockConunt, 0)
}

func scan(ctx context.Context) (*route.Model, bool) {
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	return route.Scan(ctx, Model, Meta)
}

func search(ctx context.Context, template *route.Model) (*route.Payload, bool) {
	MuRoute.Lock()
	defer MuRoute.Unlock()

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
