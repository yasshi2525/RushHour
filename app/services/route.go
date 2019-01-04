package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"
)

type ctxkey int

const (
	workerid ctxkey = iota
	taskkey
)

var cancelChannel chan string
var routingContext context.Context
var routingCancel context.CancelFunc

// IsSearching represents whether searching is executed or not.
var IsSearching bool

// IsValidRoute represents whether routing ended successfully or not.
var IsValidRoute bool

type routeResult struct {
	Status    bool
	Results   map[uint][]*entities.Node
	Total     int
	Processed int
}

type routeWorker struct {
	Name     string
	Tasks    []*entities.Company
	Cancel   context.CancelFunc
	ResultCh chan *routeResult
}

// StartRouting start searching.
func StartRouting(msg string) {
	IsSearching = true
	IsValidRoute = false

	routingContext, routingCancel = context.WithCancel(context.Background())
	searchCtx, searchCancel := context.WithCancel(routingContext)
	reflectCtx, reflectCancel := context.WithCancel(routingContext)

	go func() {
		MuRoute.Lock()
		defer MuRoute.Unlock()
		start := time.Now()

		if ok := search(searchCtx, searchCancel, msg); ok {
			if reflectTo(reflectCtx, reflectCancel, msg) {
				WarnLongExec(start, Config.Perf.Routing.D, "routing")
			}
		}
	}()
}

// CancelRouting stop current executing searching.
func CancelRouting(msg string) {
	if routingCancel != nil {
		routingCancel()
		//revel.AppLog.Infof("canceleld by %s", msg)
	}
}

func search(root context.Context, cancel context.CancelFunc, msg string) bool {
	defer cancel()
	MuStatic.RLock()
	defer MuStatic.RUnlock()

	MuDynamic.RLock()
	defer MuDynamic.RUnlock()

	parallel := Config.Routing.Worker
	collectCh := make(chan *routeResult, parallel)
	done := make(chan *routeResult)
	defer close(done)

	contexts, workers := genWorkers(root, parallel, collectCh)
	workWg := &sync.WaitGroup{}

	for i, w := range workers {
		workWg.Add(1)
		go subsearch(contexts[i], w, workWg)
	}
	go reduceResult(collectCh, done)

	workWg.Wait()
	close(collectCh)

	finalResult := <-done
	RouteTemplate = finalResult.Results
	IsValidRoute = finalResult.Status

	if finalResult.Status {
		//revel.AppLog.Debugf("routing was successfully ended (%d/%d)", finalResult.Processed, finalResult.Total)
	} else {
		revel.AppLog.Debugf("routing was canceled (%d/%d)", finalResult.Processed, finalResult.Total)
	}
	IsSearching = false
	return IsValidRoute
}

func genWorkers(root context.Context, parallel int, collectCh chan *routeResult) ([]context.Context, []*routeWorker) {
	workers := make([]*routeWorker, parallel, parallel)
	contexts := make([]context.Context, parallel, parallel)
	for i := 0; i < parallel; i++ {
		workers[i] = &routeWorker{
			Name:     fmt.Sprintf("worker%d", i+1),
			Tasks:    []*entities.Company{},
			ResultCh: collectCh,
		}
		contexts[i], workers[i].Cancel = context.WithCancel(root)
	}
	// group up to <parallel> routeWorker
	grpid := 0
	for _, c := range Model.Companies {
		workers[grpid].Tasks = append(workers[grpid].Tasks, c)
		grpid++
		if grpid == parallel {
			grpid = 0
		}
	}
	return contexts, workers
}

// subsearch is worker's task
func subsearch(ctx context.Context, w *routeWorker, workWg *sync.WaitGroup) {
	defer w.Cancel()
	defer workWg.Done()
	result := &routeResult{
		Status:  true,
		Results: make(map[uint][]*entities.Node),
		Total:   len(w.Tasks),
	}
	for _, c := range w.Tasks {
		select {
		case <-ctx.Done():
			// Done() called multiple times
			result.Status = result.Processed == result.Total
			break
		default:
			goal, nodes := genNodes(c)
			entities.GenEdges(nodes, Model.Steps)
			goal.WalkThrough()
			for _, n := range nodes {
				n.Fix()
			}
			result.Results[c.ID] = nodes
			result.Processed++
		}
	}
	w.ResultCh <- result
}

func reduceResult(collectCh chan *routeResult, done chan *routeResult) {
	result := &routeResult{
		Status:  true,
		Results: make(map[uint][]*entities.Node),
	}
	for res := range collectCh {
		result.Total += res.Total
		result.Processed += res.Processed
		if !res.Status {
			result.Status = false
		} else {
			for cid, nodes := range res.Results {
				result.Results[cid] = nodes
			}
		}
	}
	done <- result
}

// genNodes returns it's wrapper Node and all Node
func genNodes(goal *entities.Company) (*entities.Node, []*entities.Node) {
	var wrapper *entities.Node
	ns := []*entities.Node{}

	for _, res := range Meta.List {
		if _, ok := res.Obj().(entities.Relayable); ok {
			ForeachModel(res, func(obj interface{}) {
				if h, isHuman := obj.(entities.Human); isHuman && h.To != goal {
					revel.AppLog.Debugf("skip %s(%d) because dept=%s(%d)", res, h.ID, entities.COMPANY, goal.ID)
					return
				}
				n := entities.NewNode(obj.(entities.Relayable))
				if obj == goal {
					wrapper = n
					//revel.AppLog.Debugf("found wrapper %s(%d) for %s(%d)", res, obj.(entities.Indexable).Idx(), entities.COMPANY, goal.ID)
				}
				ns = append(ns, n)
			})
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
	for _, h := range Model.Humans {
		select {
		case <-ctx.Done():
			//revel.AppLog.Debugf("reflect Canceleld %s in %d / %d", msg, cnt+1, len(Model.Humans))
			return false
		default:
			for _, n := range RouteTemplate[h.To.ID] {
				if n.Base == h {
					Model.Agents[h.ID].Current = n.ViaEdge
				}
			}
			cnt++
		}
	}
	return true
}
