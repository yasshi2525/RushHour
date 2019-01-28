package entities

import (
	"fmt"

	"github.com/revel/revel"
)

// Depart creates new LineTask which deprats from it.
// if it is already connected, it panics when force is false.
func (lt *LineTask) Depart(force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to depart from Connectted LineTask: %v", lt))
	}
	if lt.TaskType != OnStopping {
		panic(fmt.Errorf("Tried to depart from invald TaskType : %v", lt))
	}
	return lt.M.NewLineTaskDept(lt.RailLine, lt.Dest, lt)
}

// DepartIf creates new LineTask which departs from it if can.
// if it is already connected, it panics when force is false.
func (lt *LineTask) DepartIf(force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to depart from Connectted LineTask: %v", lt))
	}
	if lt.TaskType == OnStopping {
		if lt.Dest == nil {
			revel.AppLog.Errorf("dest is nil %v", lt)
		}
		return lt.Depart(force...)
	}
	return lt
}

// Stretch creates new LineTask which moves on specified RailEdge.
// if it is already connected, it panics when force is false.
func (lt *LineTask) Stretch(re *RailEdge, force ...bool) *LineTask {
	if !(len(force) > 0 && force[0]) && lt.next != nil {
		panic(fmt.Errorf("Tried to add RailEdge to Connectted LineTask: %v -> %v", re, lt))
	}
	if lt.ToNode() != re.FromNode {
		panic(fmt.Errorf("Tried to add far RailEdge to LineTask: %v -> %v", re, lt))
	}

	// when task is "stop", append task "departure"
	tail := lt.DepartIf(force...)
	return lt.M.NewLineTask(lt.RailLine, re, tail)
}

// InsertRailEdge corrects RailLine for new RailEdge
// RailEdge.From must be original RailNode.
// RailEdge.To   must be      new RailPoint.
//
// Before (a) ---------------> (b) -> (c)
// After  (a) -> (X) -> (a) -> (b) -> (c)
//
// RailEdge : (a) -> (X)
func (lt *LineTask) InsertRailEdge(re *RailEdge) {
	if lt.ToNode() != re.FromNode {
		panic(fmt.Errorf("Tried to insert far RailEdge to LineTask: %v -> %v", re, lt))
	}
	next := lt.Next()                     // = (b) -> (c)
	tail := lt.Stretch(re, true)          // = (a) -> (X)
	tail = tail.Stretch(re.Reverse, true) // = (X) -> (a)

	// when (X) is station and is stopped, append "dept" task after it
	if tail.TaskType == OnStopping && next != nil && next.TaskType != OnDeparture {
		tail = tail.DepartIf()
	}
	tail.SetNext(next) // (a) -> (b) -> (c)
}

// InsertDestination set specified Platform to it's destination.
func (lt *LineTask) InsertDestination(p *Platform) {
	if lt.TaskType == OnDeparture {
		panic(fmt.Errorf("try to insert destination to dept LineTask: %v -> %v", p, lt))
	}
	lt.Dest = p
	lt.DestID = p.ID
	if lt.RailLine.AutoPass {
		// change move -> pass
		lt.TaskType = OnPassing
		lt.RailLine.ReRouting = true
	} else {
		// change move -> stop
		lt.TaskType = OnStopping
		next := lt.next
		lt.Depart(true).SetNext(next)
	}
}

// InsertDeparture set specified Platform to it's departure.
func (lt *LineTask) InsertDeparture(p *Platform) {
	lt.SetDept(p)
}

// Shrink unregister specified Platform as stop or pass target.
func (lt *LineTask) Shrink(p *Platform) {
	if lt.Stay != p {
		panic(fmt.Errorf("try to shrink far platform: %v -> %v", p, lt))
	}
	next := lt.next
	if next != nil {
		next.SetDept(nil)
	}
	lt.SetNext(nil)
	if lt.before != nil {
		lt.before.SetDest(nil)
		lt.before.TaskType = OnMoving
		lt.before.SetNext(next)
	}
	lt.Delete(false)
}

// Shave shorten it's route as skip specified RailEdge.
func (lt *LineTask) Shave(re *RailEdge) {
	if lt.Moving != re {
		panic(fmt.Errorf("try to shave far edge: %v -> %v", re, lt))
	}
	if lt.next != nil {
		if lt.next.Moving != re.Reverse {
			panic(fmt.Errorf("try to shave linear RailLine: %v -> %v", re.Reverse, lt.next))
		}
		if lt.before != nil {
			if lt.next.next != nil && lt.next.next.TaskType == OnDeparture {
				// skip redundant Departure
				if lt.before.TaskType == OnDeparture {
					lt.before.SetNext(lt.next.next.next)
					lt.next.next.Delete(false)
				} else {
					lt.before.SetNext(lt.next.next)
				}
				if lt.before.TaskType == OnPassing {
					lt.before.TaskType = OnStopping
					lt.before.SetDest(lt.next.Stay)
				}
			}

		}
		lt.next.Delete(false)
	}
	lt.Delete(false)
}
