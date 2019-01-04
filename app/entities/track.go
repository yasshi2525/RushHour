package entities

import "fmt"

// Track represents two Travelable is connected by rail.
type Track struct {
	ID   uint
	from Travelable
	to   Travelable
	cost float64
}

// NewTrack create new instance and relation to Travelable
func NewTrack(id uint, f Travelable, t Travelable, weight float64) *Track {
	track := &Track{
		ID:   id,
		from: f,
		to:   t,
		cost: f.Pos().Dist(t) * weight,
	}
	track.Init()
	f.OutTrack()[track.ID] = track
	t.InTrack()[track.ID] = track
	return track
}

// Idx returns unique id field.
func (tr *Track) Idx() uint {
	return tr.ID
}

// Type returns type of entitiy
func (tr *Track) Type() ModelType {
	return TRACK
}

// Init do nothing
func (tr *Track) Init() {
}

// From returns where Track comes from
func (tr *Track) From() Indexable {
	return tr.from
}

// To returns where Track goes to
func (tr *Track) To() Indexable {
	return tr.to
}

// Cost is calculated by distance
func (tr *Track) Cost() float64 {
	return tr.cost
}

// UnRef delete selt from related Travelable.
func (tr *Track) UnRef() {
	delete(tr.from.OutTrack(), tr.ID)
	delete(tr.to.InTrack(), tr.ID)
}

// String represents status
func (tr *Track) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v", Meta.Attr[tr.Type()].Short,
		tr.ID, tr.from, tr.to)
}
