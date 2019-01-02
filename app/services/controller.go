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

	for idx, res := range Meta.StaticList {
		list := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(res.Type())), 0, 0)
		mapdata := Meta.StaticMap[res]
		for _, key := range mapdata.MapKeys() {
			obj := mapdata.MapIndex(key).Interface()
			if pos, ok := obj.(entities.Locationable); ok {
				if pos.IsIn(center, scale) {
					list = reflect.Append(list, reflect.ValueOf(pos))
				}
			}
		}
		view.Elem().Field(idx).Set(list)
	}
	return view.Elem().Interface()
}

func newGameView() reflect.Value {
	fields := []reflect.StructField{}
	for _, res := range Meta.StaticList {
		fields = append(fields, reflect.StructField{
			Name: res.String(),
			Type: reflect.SliceOf(reflect.PtrTo(res.Type())),
			Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", Meta.Static[res].API)),
		})
	}
	return reflect.New(reflect.StructOf(fields))
}
