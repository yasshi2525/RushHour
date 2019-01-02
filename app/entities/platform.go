package entities

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Model
	Owner

	out map[uint]*Step
	in  map[uint]*Step

	InStation  *Station        `gorm:"-" json:"-"`
	OnRailNode *RailNode       `gorm:"-" json:"-"`
	Passenger  map[uint]*Human `gorm:"-" json:"-"`

	Capacity uint `gorm:"not null"`
	Occupied uint `gorm:"not null"`

	InStationID  uint `gorm:"not null"`
	OnRailNodeID uint `gorm:"not null"`
}

// Idx returns unique id field.
func (p *Platform) Idx() uint {
	return p.ID
}

// Init creates map.
func (p *Platform) Init() {
	p.Model.Init()
	p.Owner.Init()
	p.out = make(map[uint]*Step)
	p.in = make(map[uint]*Step)
	p.Passenger = make(map[uint]*Human)
}

// Pos returns location
func (p *Platform) Pos() *Point {
	return p.OnRailNode.Pos()
}

// IsIn returns it should be view or not.
func (p *Platform) IsIn(center *Point, scale float64) bool {
	return p.Pos().IsIn(center, scale)
}

// Out returns where it can go to
func (p *Platform) Out() map[uint]*Step {
	return p.out
}

// In returns where it comes from
func (p *Platform) In() map[uint]*Step {
	return p.in
}

// Resolve set reference
func (p *Platform) Resolve(rn *RailNode, st *Station) {
	p.Owner, p.OnRailNode, p.InStation = rn.Owner, rn, st
	p.ResolveRef()
}

// ResolveRef set id from reference
func (p *Platform) ResolveRef() {
	p.Owner.ResolveRef()
	p.OnRailNodeID = p.OnRailNode.ID
	p.InStationID = p.InStation.ID
}

// Permits represents Player is permitted to control
func (p *Platform) Permits(o *Player) bool {
	return p.Owner.Permits(o)
}
