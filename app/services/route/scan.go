package route

import (
	"context"

	"github.com/yasshi2525/RushHour/app/entities"
)

func Scan(ctx context.Context, model *entities.Model) (*Model, bool) {
	result := NewModel()
	result.GoalIDs = model.Ids(entities.COMPANY)

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
