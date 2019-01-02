package services

import (
	"context"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

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
		start := time.Now()
		MuRoute.Lock()
		defer MuRoute.Unlock()
		defer routingCancel()

		if ok := search(searchCtx, searchCancel, msg); ok {
			if reflectTo(reflectCtx, reflectCancel, msg) {
				WarnLongExec(start, 10, "経路探索", true)
			}
		}
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

func search(ctx context.Context, cancel context.CancelFunc, msg string) bool {
	defer cancel()
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	cnt := 0
	for _, c := range Repo.Static.Companies {
		select {
		case <-ctx.Done():
			revel.AppLog.Debugf("search Canceleld %s in %d / %d", msg, cnt+1, len(Repo.Static.Companies))
			return false
		default:
			goal, nodes := genNodes(c)
			entities.GenEdges(nodes, Repo.Dynamic.Steps)
			goal.WalkThrough()
			for _, n := range nodes {
				n.Fix()
			}
			RouteTemplate[c.ID] = nodes
			cnt++
		}
	}
	return true
}

// genNodes returns it's wrapper Node and all Node
func genNodes(goal *entities.Company) (*entities.Node, []*entities.Node) {
	var wrapper *entities.Node
	ns := []*entities.Node{}

	for _, res := range Repo.Meta.StaticList {
		if _, ok := res.Obj().(entities.Relayable); ok {
			mapdata := Repo.Meta.StaticMap[res]
			for _, key := range mapdata.MapKeys() {
				obj := mapdata.MapIndex(key).Interface()

				if h, isHuman := obj.(entities.Human); isHuman && h.To != goal {
					revel.AppLog.Debugf("skip %s(%d) because dept=%s(%d)", res, h.ID, entities.COMPANY, goal.ID)
					continue
				}

				n := entities.NewNode(obj.(entities.Relayable))
				if obj == goal {
					wrapper = n
					//revel.AppLog.Debugf("found wrapper %s(%d) for %s(%d)", res, obj.(entities.Indexable).Idx(), entities.COMPANY, goal.ID)
				}
				ns = append(ns, n)
			}
			//revel.AppLog.Debugf("gen node %s for routing len(%d)", res, len(ns))
		}
	}
	//revel.AppLog.Debugf("gen %d nodes towards %s(%d))", len(ns), entities.COMPANY, goal.ID)
	return wrapper, ns
}

func reflectTo(ctx context.Context, cancel context.CancelFunc, msg string) bool {
	defer cancel()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	cnt := 0
	for _, h := range Repo.Static.Humans {
		select {
		case <-ctx.Done():
			revel.AppLog.Debugf("reflect Canceleld %s in %d / %d", msg, cnt+1, len(Repo.Static.Humans))
			return false
		default:
			for _, n := range RouteTemplate[h.To.ID] {
				if n.Base == h {
					h.Current = n.ViaEdge
				}
			}
			cnt++
		}
	}
	return true
}
