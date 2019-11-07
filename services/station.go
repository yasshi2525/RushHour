package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/route"
)

//CreateStation create Station
func CreateStation(o *entities.Player, rn *entities.RailNode, name string) (*entities.Station, error) {
	if err := CheckAuth(o, rn); err != nil {
		return nil, err
	}
	if rn.OverPlatform != nil {
		return nil, fmt.Errorf("staiton already exists")
	}

	st := Model.NewStation(o)
	g := Model.NewGate(st)
	p := Model.NewPlatform(rn, g)

	st.Name = name
	StartRouting()
	AddOpLog("CreateStation", o, rn, st, g, p)
	return st, nil
}

//RemoveStation remove Station
func RemoveStation(o *entities.Player, id uint) error {
	if st, err := Model.DeleteIf(o, entities.STATION, id); err != nil {
		return err
	} else {
		st := st.(*entities.Station)
		if o.ReRouting {
			route.RefreshTracks(o, conf.Game.Service.Routing.Worker)
		}
		for _, l := range o.RailLines {
			if l.ReRouting {
				route.RefreshTransports(l, conf.Game.Service.Routing.Worker)
			}
		}
		StartRouting()
		AddOpLog("RemoveStation", o, st)
		return nil
	}
}
