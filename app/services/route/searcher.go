package route

import (
	"context"
	"fmt"
	"sync"

	"github.com/yasshi2525/RushHour/app/entities"
)

type Searcher struct {
	Name  string
	Tasks []uint
	Model *Model
	Ch    chan *Payload
}

func (s *Searcher) Search(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	payload := &Payload{make(map[uint]*Model), 0, len(s.Tasks)}
	for _, cid := range s.Tasks {
		select {
		case <-ctx.Done():
			break
		default:
			model, goal := s.Model.ExportWith(entities.COMPANY, cid)
			goal.WalkThrough()
			for _, n := range model.Nodes {
				n.Fix()
			}
			payload.Route[cid] = model
			payload.Processed++
		}
	}
	s.Ch <- payload
}

func Search(ctx context.Context, parallel int, model *Model) (*Payload, bool) {
	reduceCh := make(chan *Payload)
	defer close(reduceCh)

	// create worker
	subctxs, searchers, collectCh := genSearchers(ctx, parallel, model)
	wg := &sync.WaitGroup{}

	// fire task
	for i, searcher := range searchers {
		wg.Add(1)
		go searcher.Search(subctxs[i], wg)
	}
	go reduceSearch(parallel, collectCh, reduceCh)

	// join
	wg.Wait()
	close(collectCh)

	// reduce
	result := <-reduceCh
	return result, result.IsOK()
}

func genSearchers(ctx context.Context, parallel int, model *Model) ([]context.Context, []*Searcher, chan *Payload) {
	searchers := make([]*Searcher, parallel)
	subctxs := make([]context.Context, parallel)
	ch := make(chan *Payload, parallel)

	for i := 0; i < parallel; i++ {
		name := fmt.Sprintf("searcher%d", i+1)
		searchers[i] = &Searcher{name, []uint{}, model.Export(), ch}
		subctxs[i], _ = context.WithCancel(ctx)
	}

	// slice to <parallel> group
	grpid := 0
	for _, cid := range model.GoalIDs {
		searchers[grpid].Tasks = append(searchers[grpid].Tasks, cid)
		grpid++
		if grpid == parallel {
			grpid = 0
		}
	}

	return subctxs, searchers, ch
}

func reduceSearch(parallel int, collectCh chan *Payload, reduceCh chan *Payload) {
	payload := &Payload{make(map[uint]*Model), 0, parallel}
	for sub := range collectCh {
		if sub.IsOK() {
			payload.Import(sub)
		}
	}
	reduceCh <- payload
}
