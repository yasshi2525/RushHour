package services

import (
	"fmt"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreateRailLine create RailLine
func CreateRailLine(o *entities.Player, name string) (*entities.RailLine, error) {
	l := entities.NewRailLine(GenID(entities.LINE), o)
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
		if lt.Own == o && lt.To().Pos().SameAt(re.FromNode) { // = (a) -> (b)
			bases = append(bases, lt)
		}
	}

	for _, base := range bases {
		next := base.Next // = (b) -> (c)

		inter, _ := AttachLineTask(o, base, re, pass...)         // = (a) -> (X)
		inter, _ = AttachLineTask(o, inter, re.Reverse, pass...) // = (X) -> (a)

		// when (X) is station and is stopped, append "dept" task after it
		if inter.TaskType == entities.OnStopping && next != nil && next.TaskType != entities.OnDeparture {
			inter = entities.NewLineTaskDept(GenID(entities.LINETASK), inter.RailLine, inter.Dest, inter)
			AddEntity(inter)
		}

		inter.Next = next // (a) -> (b) -> (c)
	}
	return nil
}

// AttachLineTask attaches new RailEdge
func AttachLineTask(o *entities.Player, tail *entities.LineTask, newer *entities.RailEdge, pass ...bool) (*entities.LineTask, error) {
	if err := CheckAuth(o, tail); err != nil {
		return nil, err
	}
	if err := CheckAuth(o, newer); err != nil {
		return nil, err
	}
	if !tail.To().Pos().SameAt(newer.From()) {
		return nil, fmt.Errorf("unconnected RailEdge. %v != %v", tail.To().Pos(), newer.From().Pos())
	}

	// when task is "stop", append task "departure"
	if tail.TaskType == entities.OnStopping {
		tail = entities.NewLineTaskDept(GenID(entities.LINETASK), tail.RailLine, tail.Stay, tail)
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
	if head, tail, _ := FindRailLineBorder(o, l); head != nil && tail != nil {
		tail.Next = head
	}
	return false, nil
}

// FindRailLineBorder returns head and tail of LineTask.
// Head and tail are nil when LineTask loops
// Tail is undirecting LineTask, that is LineTask.Next is nil
// Head is undirected  LineTask because head of chain is what any other doesn't target
func FindRailLineBorder(o *entities.Player, l *entities.RailLine) (*entities.LineTask, *entities.LineTask, error) {
	if err := CheckAuth(o, l); err != nil {
		return nil, nil, err
	}
	var tail *entities.LineTask

	referred := make(map[uint]bool)
	for _, lt := range l.Tasks {
		if lt.Next != nil {
			referred[lt.Next.ID] = true
		} else {
			tail = lt
		}
	}

	for key, v := range referred {
		if !v {
			return l.Tasks[key], tail, nil
		}
	}
	// looped
	return nil, nil, nil
}
