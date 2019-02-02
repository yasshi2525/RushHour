package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
)

// ViewMap immitates user requests view
func ViewMap(x float64, y float64, scale float64, after ...time.Time) interface{} {
	start := time.Now()
	MuModel.RLock()
	defer MuModel.RUnlock()
	lock := time.Now()
	defer WarnLongExec(start, lock, Const.Perf.View.D, "view")

	view := newGameView()

	for idx, res := range entities.TypeList {
		// filter out agent, step ...
		if res.IsVisible() {
			list := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(res.Type())), 0, 0)
			Model.ForEach(res, func(obj entities.Entity) {
				// filter time
				if len(after) > 0 && !obj.B().IsChanged(after[0]) {
					return
				}
				// filter out of user view
				if obj.S().IsIn(x, y, scale) {
					list = reflect.Append(list, reflect.ValueOf(obj))
				}
			})
			view.Elem().Field(idx).Set(list)
		}
	}
	return view.Elem().Interface()
}

func newGameView() reflect.Value {
	fields := []reflect.StructField{}
	for _, res := range entities.TypeList {
		// filter out agent, step ...
		if res.IsVisible() {
			fields = append(fields, reflect.StructField{
				Name: res.String(),
				Type: reflect.SliceOf(reflect.PtrTo(res.Type())),
				Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", res.API())),
			})
		}
	}
	return reflect.New(reflect.StructOf(fields))
}

// CheckAuth throws error when there is no permission
func CheckAuth(owner *entities.Player, res entities.Entity) error {
	if res.B().Permits(owner) {
		return nil
	}
	return fmt.Errorf("no permission to operate %v", res)
}
