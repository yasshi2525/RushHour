package entities

type DelegateRailNode struct {
	Base
	Point

	RailNodes map[uint]*RailNode `json:"-"`
}

type DelegateRailEdge struct {
	Base

	From *DelegateRailNode `json:"-"`
	To   *DelegateRailNode `json:"-"`

	Tracks map[uint]*Track `json:"-"`

	FromID    uint `json:"from"`
	ToID      uint `json:"to"`
	ReverseID uint `json:"eid"`
}
