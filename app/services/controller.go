package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
)

// ViewMap immitates user requests view
func ViewMap(x float64, y float64, scale float64) interface{} {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	center := &entities.Point{X: x, Y: y}
	view := newGameView()

	for idx, res := range Meta.List {
		// filter agent, step ...
		if res.IsVisible() {
			list := reflect.MakeSlice(reflect.SliceOf(res.Type()), 0, 0)
			ForeachModel(res, func(obj interface{}) {
				// filter out of user view
				if pos, ok := obj.(entities.Locationable); ok && pos.IsIn(center, scale) {
					list = reflect.Append(list, reflect.ValueOf(pos))
				}
			})
			view.Elem().Field(idx).Set(list)
		}
	}
	return view.Elem().Interface()
}

func newGameView() reflect.Value {
	fields := []reflect.StructField{}
	for _, res := range Meta.List {
		fields = append(fields, reflect.StructField{
			Name: res.String(),
			Type: reflect.SliceOf(res.Type()),
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", res.API())),
		})
	}
	return reflect.New(reflect.StructOf(fields))
}
