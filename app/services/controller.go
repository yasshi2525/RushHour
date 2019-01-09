package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// ViewMap immitates user requests view
func ViewMap(x float64, y float64, scale float64, after ...time.Time) interface{} {
	start := time.Now()
	defer WarnLongExec(start, Config.Perf.View.D, "view")

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	view := newGameView()

	for idx, res := range entities.TypeList {
		// filter out agent, step ...
		if res.IsVisible() {
			list := reflect.MakeSlice(reflect.SliceOf(res.Type()), 0, 0)
			Model.ForEach(res, func(obj entities.Indexable) {
				// filter time
				if len(after) > 0 && !obj.(entities.Persistable).IsChanged(after[0]) {
					return
				}
				// filter out of user view
				if pos, ok := obj.(entities.Viewable); ok && pos.IsIn(x, y, scale) {
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
	for _, res := range entities.TypeList {
		// filter out agent, step ...
		if res.IsVisible() {
			fields = append(fields, reflect.StructField{
				Name: res.String(),
				Type: reflect.SliceOf(res.Type()),
				Tag:  reflect.StructTag(fmt.Sprintf("json:\"%s\"", res.API())),
			})
		}
	}
	return reflect.New(reflect.StructOf(fields))
}

// CheckAuth throws error when there is no permission
func CheckAuth(owner *entities.Player, res entities.Ownable) error {
	if res.Permits(owner) {
		return nil
	}
	return fmt.Errorf("no permission to operate %v", res)
}

// TryRemove delete entity if it can.
func TryRemove(
	o *entities.Player,
	res entities.ModelType,
	id uint,
	callback func(interface{})) error {

	v := Model.Values[res].MapIndex(reflect.ValueOf(id))

	// no id
	if !v.IsValid() {
		revel.AppLog.Infof("%v(%d) was already removed.", res, id)
		return nil
	}

	obj := v.Interface()

	// no permission
	if prop, ok := obj.(entities.Ownable); ok && !prop.Permits(o) {
		revel.AppLog.Infof("%v is not permitted to delete %v", o, prop)
		return fmt.Errorf("no permission to delete %v", prop)
	}

	// reference
	if err := obj.(entities.Removable).CheckRemove(); err != nil {
		return err
	}

	callback(obj)
	return nil
}
