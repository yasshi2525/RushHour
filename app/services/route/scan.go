package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func Scan(ctx context.Context, model *entities.Model, meta *entities.MetaModel) (*Model, bool) {
	result := NewModel()
	result.GoalIDs = model.Ids(entities.COMPANY)

	if !genNodes(ctx, result, model, meta) {
		return result, false
	}

	for _, s := range model.Steps {
		select {
		case <-ctx.Done():
			return result, false
		default:
			result.FindOrCreateEdge(s)
		}
	}
	return result, true
}

func genNodes(ctx context.Context, result *Model, model *entities.Model, meta *entities.MetaModel) bool {
	for _, res := range meta.List {
		if _, ok := res.Obj().(entities.Relayable); ok {
			select {
			case <-ctx.Done():
				return false
			default:
				mapdata := meta.Map[res]
				for _, key := range mapdata.MapKeys() {
					obj := mapdata.MapIndex(key).Interface().(entities.Indexable)
					result.FindOrCreateNode(obj)
				}
			}
		}
	}
	return true
}
