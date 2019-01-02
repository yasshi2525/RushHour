package services

import (
	"fmt"
	"reflect"
	"time"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/entities"
)

// Restore get model from database
func Restore() {
	revel.AppLog.Info("start restore from database")
	defer revel.AppLog.Info("end restore from database")

	start := time.Now()
	defer WarnLongExec(start, 5, "DBリストア", true)

	MuStatic.Lock()
	defer MuStatic.Unlock()
	MuDynamic.Lock()
	defer MuDynamic.Unlock()

	setNextID()
	fetchStatic()
	resolveStatic()
	genDynamics()
}

// setNextID set max id as NextID from database for Restore()
func setNextID() {
	for _, key := range Meta.StaticList {
		var maxID struct {
			V uint64
		}
		sql := fmt.Sprintf("SELECT max(id) as v FROM %s", key.Table())
		if err := db.Raw(sql).Scan(&maxID).Error; err == nil {
			Static.NextIDs[key] = &maxID.V
			revel.AppLog.Debugf("set NextID[%s] = %d", key, *Static.NextIDs[key])
		} else {
			panic(err)
		}
	}
}

// fetchStatic selects records for Restore()
func fetchStatic() {
	for _, key := range Meta.StaticList {
		// select文組み立て
		if rows, err := db.Table(key.Table()).Where("deleted_at is null").Rows(); err == nil {
			for rows.Next() {
				// 対応する Struct を作成
				base := key.Obj()
				if err := db.ScanRows(rows, base); err == nil {
					// Static に登録
					if obj, ok := base.(entities.Indexable); ok {
						Meta.StaticMap[key].SetMapIndex(reflect.ValueOf(obj.Idx()), reflect.ValueOf(obj))
						//revel.AppLog.Debugf("set Static[%s][%d] = %+v", key, obj.Idx(), obj)
					} else {
						panic(fmt.Errorf("invalid type %T: %+v", base, base))
					}
				} else {
					panic(err)
				}
			}
		} else {
			panic(err)
		}

	}
}

// resolveStatic set pointer from id for Restore()
func resolveStatic() {
	for _, rn := range Static.RailNodes {
		rn.Resolve(Static.Players[rn.OwnerID])
	}
	for _, re := range Static.RailEdges {
		re.Resolve(Static.RailNodes[re.FromID], Static.RailNodes[re.ToID])
	}
	for _, st := range Static.Stations {
		st.Resolve(Static.Players[st.OwnerID])
	}
	for _, g := range Static.Gates {
		g.Resolve(Static.Stations[g.InStationID])
	}
	for _, p := range Static.Platforms {
		p.Resolve(Static.RailNodes[p.OnRailNodeID], Static.Stations[p.InStationID])
	}
	for _, l := range Static.RailLines {
		l.Resolve(Static.Players[l.OwnerID])
	}
	for _, lt := range Static.LineTasks {
		lt.Resolve(Static.RailLines[lt.RailLineID])
		// nullable fields
		if lt.NextID != 0 {
			lt.Resolve(Static.LineTasks[lt.NextID])
		}
		if lt.StayID != 0 {
			lt.Resolve(Static.Platforms[lt.StayID])
		}
		if lt.MovingID != 0 {
			lt.Resolve(Static.RailEdges[lt.MovingID])
		}
	}
	for _, t := range Static.Trains {
		t.Resolve(Static.LineTasks[t.TaskID])
	}
	for _, h := range Static.Humans {
		h.Resolve(Static.Residences[h.FromID], Static.Companies[h.ToID])
		// nullable fields
		if h.OnPlatformID != 0 {
			h.Resolve(Static.Platforms[h.OnPlatformID])
		}
		if h.OnTrainID != 0 {
			h.Resolve(Static.Platforms[h.OnTrainID])
		}
	}
}

// genDynamics create Dynamic instances
func genDynamics() {
	walk, train := Config.Human.Weight, Config.Train.Weight
	for _, r := range Static.Residences {
		// R -> C, G
		GenStepResidence(r)
	}
	for _, c := range Static.Companies {
		// G -> C
		for _, g := range Static.Gates {
			GenStep(g, c, walk)
		}
	}
	for _, p := range Static.Platforms {
		// G <-> P
		g := p.InStation.Gate
		GenStep(p, g, walk)
		GenStep(g, p, walk)

		// P <-> P
		for _, p2 := range Static.Platforms {
			if p != p2 {
				GenStep(p, p2, train)
			}
		}
	}
	for _, h := range Static.Humans {
		GenStepHuman(h)
		Dynamic.Agents[h.ID] = entities.NewAgent(h)
	}
}
