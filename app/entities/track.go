package entities

import "fmt"

type Track struct {
	ID       uint
	M        *Model
	FromNode *RailNode
	ToNode   *RailNode
	Via      *RailEdge
	Value    float64
}

func (m *Model) NewTrack(f *RailNode, t *RailNode, via *RailEdge, v float64) *Track {
	tk := &Track{
		ID:       m.GenID(TRACK),
		FromNode: f,
		ToNode:   t,
		Via:      via,
		Value:    v,
	}
	tk.Init(m)
	f.Tracks[t.ID] = tk
	m.Add(tk)
	return tk
}

// Idx returns unique id field.
func (tk *Track) Idx() uint {
	return tk.ID
}

// Type returns type of entitiy
func (tk *Track) Type() ModelType {
	return TRACK
}

// Init do nothing
func (tk *Track) Init(m *Model) {
	tk.M = m
}

// From returns where Track comes from
func (tk *Track) From() Indexable {
	return tk.FromNode
}

// To returns where Track goes to
func (tk *Track) To() Indexable {
	return tk.ToNode
}

// Cost is calculated by distance
func (tk *Track) Cost() float64 {
	return tk.Value
}

// Pos returns center
func (tk *Track) Pos() *Point {
	return tk.Via.Pos()
}

func (tk *Track) IsIn(x float64, y float64, scale float64) bool {
	return tk.Via.IsIn(x, y, scale)
}

// String represents status
func (tk *Track) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,via=%v,v=%.2f", tk.Type().Short(),
		tk.ID, tk.FromNode, tk.ToNode, tk.Via, tk.Cost())
}
