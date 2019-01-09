package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func Scan(ctx context.Context, model *entities.Model) (*Model, bool) {
	result := NewModel()
	result.GoalIDs = model.Ids(entities.COMPANY)

	if !genNodes(ctx, result, model) {
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

func genNodes(ctx context.Context, result *Model, model *entities.Model) bool {
	for _, res := range entities.TypeList {
		if _, ok := res.Obj().(entities.Relayable); ok {
			select {
			case <-ctx.Done():
				return false
			default:
				model.ForEach(res, func(obj entities.Indexable) {
					result.FindOrCreateNode(obj)
				})
			}
		}
	}
	return true
}
