package entities

import "fmt"

// Transport has minimum distance route on line which connects two Platform.
type Transport struct {
	Base
	Shape
	FromPlatform *Platform
	ToPlatform   *Platform
	Via          *LineTask
	Value        float64
}

// NewTransport creates instance.
func (m *Model) NewTransport(f *Platform, t *Platform, via *LineTask, v float64) *Transport {
	x := &Transport{
		Base:         m.NewBase(TRANSPORT),
		Shape:        NewShapeEdge(&f.OnRailNode.Point, &t.OnRailNode.Point),
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

// B returns base information of this elements.
func (x *Transport) B() *Base {
	return &x.Base
}

// S returns entities' position.
func (x *Transport) S() *Shape {
	return &x.Shape
}

// Init do nothing
func (x *Transport) Init(m *Model) {
	x.Base.Init(TRANSPORT, m)
}

// From returns where Track comes from
func (x *Transport) From() Entity {
	return x.FromPlatform
}

// To returns where Track goes to
func (x *Transport) To() Entity {
	return x.ToPlatform
}

// Cost represents how many seconds it takes.
func (x *Transport) Cost() float64 {
	return x.Value
}

// Resolve set reference from id.
func (x *Transport) Resolve(args ...Entity) {
}

// CheckDelete checks related reference
func (x *Transport) CheckDelete() error {
	return nil
}

// BeforeDelete delete selt from related Locationable.
func (x *Transport) BeforeDelete() {
	delete(x.FromPlatform.Transports, x.ToPlatform.ID)
}

// Delete removes this entity with related ones.
func (x *Transport) Delete() {
	x.M.Delete(x)
}

// String represents status
func (x *Transport) String() string {
	return fmt.Sprintf("%s(%v):from=%v,to=%v,via=%v,v=%.2f", x.Type().Short(),
		x.ID, x.FromPlatform, x.ToPlatform, x.Via, x.Value)
}
