package services

import (
	"fmt"
	"log"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/route"
)

// CreateRailLine create RailLine
func CreateRailLine(o *entities.Player, name string, ext bool, pass bool) (*entities.RailLine, error) {
	l := Model.NewRailLine(o)
	l.Name = name
	l.AutoExt = ext
	l.AutoPass = pass

	StartRouting()
	AddOpLog("CreateRailLine", o, l)
	return l, nil
}

// StartRailLine start RailLine at Station
func StartRailLine(
	o *entities.Player,
	l *entities.RailLine,
	p *entities.Platform) error {
	if err := CheckAuth(o, l); err != nil {
		return err
	}
	if err := CheckAuth(o, p); err != nil {
		return err
	}
	if len(l.Tasks) > 0 {
		return fmt.Errorf("task is already registered: %v", l)
	}
	l.StartPlatform(p)
	if l.ReRouting {
		route.RefreshTransports(l, serviceConf.AppConf.Game.Service.Routing.Worker)
	}
	StartRouting()
	AddOpLog("StartRailLine", o, l, p)
	return nil
}

func StartRailLineEdge(o *entities.Player,
	l *entities.RailLine,
	re *entities.RailEdge) error {
	if err := CheckAuth(o, l); err != nil {
		return err
	}
	if err := CheckAuth(o, re); err != nil {
		return err
	}
	if len(l.Tasks) > 0 {
		return fmt.Errorf("task is already registered: %v", l)
	}
	l.StartEdge(re)
	if l.ReRouting {
		route.RefreshTransports(l, serviceConf.AppConf.Game.Service.Routing.Worker)
	}
	StartRouting()
	AddOpLog("StartRailLineEdge", o, l, re)
	return nil
}

func InsertLineTaskRailEdge(o *entities.Player, l *entities.RailLine, re *entities.RailEdge) error {
	if err := CheckAuth(o, re); err != nil {
		return err
	}
	for _, lt := range re.FromNode.InTasks {
		if lt.RailLine == l {
			lt.InsertRailEdge(re)
		}
	}
	if l.ReRouting {
		route.RefreshTransports(l, serviceConf.AppConf.Game.Service.Routing.Worker)
	}
	StartRouting()
	AddOpLog("InsertLineTaskRailEdge", o, l, re)
	return nil
}

func ComplementRailLine(o *entities.Player, l *entities.RailLine) (bool, error) {
	if err := CheckAuth(o, l); err != nil {
		return false, err
	}
	if len(l.Tasks) == 0 || l.IsRing() {
		return false, fmt.Errorf("line is already ringed: %v", l)
	}
	l.Complement()
	StartRouting()
	return true, nil
}

// RingRailLine connects tail and head
func RingRailLine(o *entities.Player, l *entities.RailLine) (bool, error) {
	if err := CheckAuth(o, l); err != nil {
		return false, err
	}
	// Check RainLine is not ringing
	ret := l.RingIf()
	if ret {
		route.RefreshTransports(l, serviceConf.AppConf.Game.Service.Routing.Worker)
		StartRouting()
		AddOpLog("RingRailLine", o, l)
	}
	return ret, nil
}

func RemoveRailLine(o *entities.Player, id uint) error {
	if l, err := Model.DeleteIf(o, entities.RAILLINE, id); err != nil {
		return err
	} else {
		StartRouting()
		AddOpLog("RemoveRailLine", o, l)
		return nil
	}
}

// [DEBUG]
func lineValidation(l *entities.RailLine) {
	var headCnt, tailCnt, loopSize int
	var deadloop, smallloop bool

	for _, lt := range l.Tasks {
		if lt.Before() == nil {
			headCnt++
		}
		if lt.Next() == nil {
			tailCnt++
		}
	}

	if headCnt > 1 {
		log.Printf("[DEBUG] MULTI HEAD DETECTED!")
	}

	if tailCnt > 1 {
		log.Printf("[DEBUG] MULTI TAIL DETECTED!")
	}

	var top *entities.LineTask
	for _, top = range l.Tasks {
		break
	}

	if top != nil {
		lt := top.Next()
		for lt != nil && lt != top {
			lt = lt.Next()
			if loopSize > len(l.Tasks) {
				log.Printf("[DEBUG] DEAD LOOP DETECTED: lt(%d)", lt.ID)
				deadloop = true
				break
			}
			loopSize++
		}
		if lt == top && loopSize < len(l.Tasks)-1 {
			log.Printf("[DEBUG] SMALL LOOP DETECTED: lt(%d)", lt.ID)
			smallloop = true
		}
	}

	if headCnt > 1 || tailCnt > 1 || deadloop || smallloop {
		dumpRailLine(l)
		time.Sleep(2 * time.Second)
		panic("error detected")
	}
}

func dumpRailLine(l *entities.RailLine) {
	for _, lt := range l.Tasks {
		log.Printf("[DEBUG] %v", lt)
	}
}
