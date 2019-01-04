package services

import (
	"github.com/yasshi2525/RushHour/app/entities"
)

//CreateStation create Station
func CreateStation(o *entities.Player, rn *entities.RailNode, name string) (*entities.Station, error) {

	if err := CheckAuth(o, rn); err != nil {
		return nil, err
	}

	st := entities.NewStation(GenID(entities.STATION), o)
	g := entities.NewGate(GenID(entities.GATE), st)
	p := entities.NewPlatform(GenID(entities.PLATFORM), rn, g, st)

	st.Name = name
	g.Num = Config.Gate.Num
	p.Capacity = Config.Platform.Capacity

	AddEntity(st, g, p)
	genStepStation(st)
	return st, nil
}

//RemoveStation remove Station
func RemoveStation(o *entities.Player, id uint) error {
	return TryRemove(o, entities.STATION, id, func(obj interface{}) {
		st := obj.(*entities.Station)
		delStepStation(st)
		DelEntity(st.Platform, st.Gate, st)
	})
}

// GenStepStation generates step related Station.
func genStepStation(st *entities.Station) {
	// R -> P
	for _, r := range Model.Residences {
		GenWalkStep(r, st.Platform)
	}
	// G <-> P
	GenWalkStep(st.Gate, st.Platform)
	GenWalkStep(st.Platform, st.Gate)
	// G -> C
	for _, c := range Model.Companies {
		GenWalkStep(st.Gate, c)
	}
}

func delStepStation(st *entities.Station) {
	for _, s := range st.Platform.InStep() {
		DelEntity(s)
	}
	for _, s := range st.Platform.OutStep() {
		DelEntity(s)
	}
	for _, s := range st.Gate.InStep() {
		DelEntity(s)
	}
	for _, s := range st.Gate.OutStep() {
		DelEntity(s)
	}
}
