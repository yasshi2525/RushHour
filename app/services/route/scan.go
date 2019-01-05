package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func Scan(ctx context.Context, model *entities.Model, meta *entities.MetaModel) (*Model, bool) {
	cids := make([]uint, len(model.Companies))
	var ns []*Node
	var es []*Edge
	var ok bool

	i := 0
	for id := range model.Companies {
		cids[i] = id
		i++
	}

	if ns, ok = genNodes(ctx, model, meta); !ok {
		//revel.AppLog.Debugf("genNodes canceled at %d/%d", len(ns), len(meta.List))
		return nil, false
	}
	if es, ok = genEdges(ctx, ns, model.Steps); !ok {
		//revel.AppLog.Debugf("genEdges canceled at %d/%d", len(es), len(model.Steps))
		return nil, false
	}
	return &Model{cids, ns, es}, true
}

func genNodes(ctx context.Context, model *entities.Model, meta *entities.MetaModel) ([]*Node, bool) {
	ns := []*Node{}
	for _, res := range meta.List {
		if _, ok := res.Obj().(entities.Relayable); ok {
			select {
			case <-ctx.Done():
				return ns, false
			default:
				mapdata := meta.Map[res]
				for _, key := range mapdata.MapKeys() {
					obj := mapdata.MapIndex(key).Interface()
					ns = append(ns, NewNode(obj.(entities.Indexable)))
				}
			}
		}
	}
	return ns, true
}

func genEdges(ctx context.Context, nodes []*Node, steps map[uint]*entities.Step) ([]*Edge, bool) {
	es := []*Edge{}
	for _, s := range steps {
		select {
		case <-ctx.Done():
			//revel.AppLog.Debugf("genEdges Canceleld in %d/%d", len(es), len(steps))
			return es, false
		default:
			es = append(es, NewEdgeFrom(nodes, s))
		}
	}
	return es, true
}
