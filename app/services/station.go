package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/route"
)

//CreateStation create Station
func CreateStation(o *entities.Player, rn *entities.RailNode, name string) (*entities.Station, error) {
	if err := CheckMaintenance(); err != nil {
		return nil, err
	}
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
	AddOpLog("CreateStation", o, rn, st, g, p)
	return st, nil
}

//RemoveStation remove Station
func RemoveStation(o *entities.Player, id uint) error {
	if err := CheckMaintenance(); err != nil {
		return err
	}
	if st, err := Model.DeleteIf(o, entities.STATION, id); err != nil {
		return err
	} else {
		st := st.(*entities.Station)
		if o.ReRouting {
			route.RefreshTracks(o, Const.Routing.Worker)
		}
		for _, l := range o.RailLines {
			if l.ReRouting {
				route.RefreshTransports(l, Const.Routing.Worker)
			}
		}
		AddOpLog("RemoveStation", o, st)
		return nil
	}
}
