package entities

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	Model
	Owner

	out map[uint]*Step
	in  map[uint]*Step

	InStation *Station `gorm:"-" json:"-"`

	// Num represents how many Human can pass at the same time
	Num uint `gorm:"not null"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null"`
	// Occupied represents how many Gate are used by Human.
	Occupied uint `gorm:"not null"`

	InStationID uint `gorm:"not null"`
}

// Idx returns unique id field.
func (g *Gate) Idx() uint {
	return g.ID
}

// Init creates map.
func (g *Gate) Init() {
	g.Model.Init()
	g.Owner.Init()
	g.out = make(map[uint]*Step)
	g.in = make(map[uint]*Step)
}

// Pos returns location
func (g *Gate) Pos() *Point {
	return g.InStation.Pos()
}

// IsIn returns it should be view or not.
func (g *Gate) IsIn(center *Point, scale float64) bool {
	return g.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (g *Gate) Out() map[uint]*Step {
	return g.out
}

// In returns where it comes from
func (g *Gate) In() map[uint]*Step {
	return g.in
}

// Resolve set reference
func (g *Gate) Resolve(st *Station) {
	g.Owner, g.InStation = st.Owner, st
	g.ResolveRef()
}

// ResolveRef set id from reference
func (g *Gate) ResolveRef() {
	g.Owner.ResolveRef()
	g.InStationID = g.InStation.ID
}

// Permits represents Player is permitted to control
func (g *Gate) Permits(o *Player) bool {
	return g.Owner.Permits(o)
}
