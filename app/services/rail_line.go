package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services/route"
)

// CreateRailLine create RailLine
func CreateRailLine(o *entities.Player, name string) (*entities.RailLine, error) {
	l := entities.NewRailLine(GenID(entities.RAILLINE), o, Config.Train.Slowness)
	l.Name = name
	AddEntity(l)
	return l, nil
}

// StartRailLine start RailLine at Station
func StartRailLine(
	o *entities.Player,
	l *entities.RailLine,
	p *entities.Platform) (*entities.LineTask, error) {
	if err := CheckAuth(o, l); err != nil {
		return nil, err
	}
	if err := CheckAuth(o, p); err != nil {
		return nil, err
	}
	if len(l.Tasks) > 0 {
		return nil, fmt.Errorf("already registered %v", l.Tasks)
	}
	lt := entities.NewLineTaskDept(GenID(entities.LINETASK), l, p)
	AddEntity(lt)
	return lt, nil
}

// InsertLineTask corrects RailLine for new RailEdge
// RailEdge.From must be original RailNode.
// RailEdge.To   must be      new RailPoint.
//
// Before (a) ---------------> (b) -> (c)
// After  (a) -> (X) -> (a) -> (b) -> (c)
//
// RailEdge : (a) -> (X)
func InsertLineTask(o *entities.Player, re *entities.RailEdge, pass ...bool) error {
	if err := CheckAuth(o, re); err != nil {
		return err
	}

	// extract tasks which direct origin
	// find (a) -> (b)
	bases := []*entities.LineTask{}
	for _, lt := range Model.LineTasks {
		if lt.Own == o && lt.ToLoc().Pos().SameAt(re.FromNode) { // = (a) -> (b)
			bases = append(bases, lt)
		}
	}

	for _, base := range bases {
		next := base.Next() // = (b) -> (c)

		inter, _ := AttachLineTask(o, base, re, pass...)         // = (a) -> (X)
		inter, _ = AttachLineTask(o, inter, re.Reverse, pass...) // = (X) -> (a)

		// when (X) is station and is stopped, append "dept" task after it
		if inter.TaskType == entities.OnStopping && next != nil && next.TaskType != entities.OnDeparture {
			inter = entities.NewLineTaskDept(GenID(entities.LINETASK), inter.RailLine, inter.Dest, inter)
			AddEntity(inter)
		}
		inter.SetNext(next) // (a) -> (b) -> (c)

		// recalculate transport cost if RailLine loops
		if inter.RailLine.IsRing() {
			delStepRailLine(inter.RailLine)
			genStepRailLine(inter.RailLine)
		}
	}
	return nil
}

// AttachLineTask attaches new RailEdge
// Need to update Step after call
func AttachLineTask(o *entities.Player, tail *entities.LineTask, newer *entities.RailEdge, pass ...bool) (*entities.LineTask, error) {
	if err := CheckAuth(o, tail); err != nil {
		return nil, err
	}
	if err := CheckAuth(o, newer); err != nil {
		return nil, err
	}
	if !tail.ToLoc().Pos().SameAt(newer.From()) {
		return nil, fmt.Errorf("unconnected RailEdge. %v != %v", tail.ToLoc().Pos(), newer.From().Pos())
	}

	// when task is "stop", append task "departure"
	if tail.TaskType == entities.OnStopping {
		tail = entities.NewLineTaskDept(GenID(entities.LINETASK), tail.RailLine, tail.Dest, tail)
		AddEntity(tail)
	}

	tail = entities.NewLineTask(GenID(entities.LINETASK), tail, newer, pass...)
	AddEntity(tail)

	return tail, nil
}

// RingRailLine connects tail and head
func RingRailLine(o *entities.Player, l *entities.RailLine) (bool, error) {
	if err := CheckAuth(o, l); err != nil {
		return false, err
	}
	// Check RainLine is not ringing
	if l.CanRing() {
		head, tail := l.Borders()
		tail.SetNext(head)
		genStepRailLine(l)
		return true, nil
	}
	return false, nil
}

// delStepRailLine discards old step
func delStepRailLine(l *entities.RailLine) {
	for _, s := range Model.Steps {
		if s.By != nil && s.By.RailLine == l {
			DelEntity(s)
		}
	}
}

// genStepRailLine generate Step P <-> P
func genStepRailLine(l *entities.RailLine) {
	tracks := route.SearchRailLine(l, Config.Routing.Worker)
	for _, tr := range tracks {
		AddEntity(tr.ExportStep(GenID(entities.STEP)))
	}
}
