package entities

import "fmt"

type Track struct {
	Base
	Shape
	FromNode *RailNode
	ToNode   *RailNode
	Via      *RailEdge
	Value    float64
}

func (m *Model) NewTrack(f *RailNode, t *RailNode, via *RailEdge, v float64) *Track {
	tk := &Track{
		Base:     m.NewBase(TRACK),
		Shape:    NewShapeEdge(&f.Point, &t.Point),
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

// B returns base information of this elements.
func (tk *Track) B() *Base {
	return &tk.Base
}

// S returns entities' position.
func (tk *Track) S() *Shape {
	return &tk.Shape
}

// Init do nothing
func (tk *Track) Init(m *Model) {
	tk.Base.Init(TRACK, m)
}

// From returns where Track comes from
func (tk *Track) From() Entity {
	return tk.FromNode
}

// To returns where Track goes to
func (tk *Track) To() Entity {
	return tk.ToNode
}

// Cost is calculated by distance
func (tk *Track) Cost() float64 {
	return tk.Value
}

// Resolve set reference from id.
func (tk *Track) Resolve(args ...Entity) {
}

func (tk *Track) CheckDelete() error {
	return nil
}

// BeforeDelete delete selt from related Locationable.
func (tk *Track) BeforeDelete() {
	delete(tk.FromNode.Tracks, tk.ToNode.ID)
}

func (tk *Track) Delete(force bool) {
	tk.M.Delete(tk)
}

// String represents status
func (tk *Track) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,via=%v,v=%.2f", tk.Type().Short(),
		tk.ID, tk.FromNode, tk.ToNode, tk.Via, tk.Cost())
}
