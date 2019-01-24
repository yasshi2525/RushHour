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

	return result, genEdges(ctx, result, model)
}

func genNodes(ctx context.Context, result *Model, model *entities.Model) bool {
	for _, res := range entities.TypeList {
		if _, ok := res.Obj(model).(entities.Relayable); ok {
			select {
			case <-ctx.Done():
				return false
			default:
				model.ForEach(res, func(obj entities.Entity) {
					result.FindOrCreateNode(obj)
				})
			}
		}
	}
	return true
}

func genEdges(ctx context.Context, result *Model, model *entities.Model) bool {
	for _, x := range model.Transports {
		select {
		case <-ctx.Done():
			return false
		default:
			result.FindOrCreateEdge(x)
		}
	}
	for _, s := range model.Steps {
		select {
		case <-ctx.Done():
			return false
		default:
			result.FindOrCreateEdge(s)
		}
	}
	return true
}
