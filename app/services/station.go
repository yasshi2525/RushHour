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
	walk, train := Config.Human.Weight, Config.Train.Weight
	// R -> P
	for _, r := range Model.Residences {
		GenStep(r, st.Platform, walk)
	}
	// G <-> P
	GenStep(st.Gate, st.Platform, walk)
	GenStep(st.Platform, st.Gate, walk)
	// P <-> P
	for _, p2 := range Model.Platforms {
		if st.Platform != p2 {
			GenStep(st.Platform, p2, train)
			GenStep(p2, st.Platform, train)
		}
	}
	// G -> C
	for _, c := range Model.Companies {
		GenStep(st.Gate, c, walk)
	}
}

func delStepStation(st *entities.Station) {
	for _, s := range st.Platform.In() {
		DelEntity(s)
	}
	for _, s := range st.Platform.Out() {
		DelEntity(s)
	}
	for _, s := range st.Gate.In() {
		DelEntity(s)
	}
	for _, s := range st.Gate.Out() {
		DelEntity(s)
	}
}
