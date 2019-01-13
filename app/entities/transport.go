package entities

import "fmt"

type Transport struct {
	ID           uint
	M            *Model
	FromPlatform *Platform
	ToPlatform   *Platform
	Via          *LineTask
	Value        float64
}

func (m *Model) NewTransport(f *Platform, t *Platform, via *LineTask, v float64) *Transport {
	x := &Transport{
		ID:           m.GenID(TRANSPORT),
		FromPlatform: f,
		ToPlatform:   t,
		Via:          via,
		Value:        v,
	}
	x.Init(m)
	f.Transports[t.ID] = x
	m.Add(x)
	return x
}

// Idx returns unique id field.
func (x *Transport) Idx() uint {
	return x.ID
}

// Type returns type of entitiy
func (x *Transport) Type() ModelType {
	return TRANSPORT
}

// Init do nothing
func (x *Transport) Init(m *Model) {
	x.M = m
}

// From returns where Track comes from
func (x *Transport) From() Indexable {
	return x.FromPlatform
}

// To returns where Track goes to
func (x *Transport) To() Indexable {
	return x.ToPlatform
}

func (x *Transport) Cost() float64 {
	return x.Value * Const.Train.Weight
}

// Pos returns center
func (x *Transport) Pos() *Point {
	return x.Via.Pos()
}

func (x *Transport) IsIn(xv float64, yv float64, scale float64) bool {
	return x.Via.IsIn(xv, yv, scale)
}

// String represents status
func (x *Transport) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,via=%v,v=%.2f", x.Type().Short(),
		x.ID, x.FromPlatform, x.ToPlatform, x.Via, x.Value)
}
