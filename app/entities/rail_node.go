package entities

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Model
	Owner
	Point
	InEdge       map[uint]*RailEdge `gorm:"-" json:"-"`
	OutEdge      map[uint]*RailEdge `gorm:"-" json:"-"`
	OverPlatform *Platform          `gorm:"-" json:"-"`
}

// NewRailNode create new instance.
func NewRailNode(id uint, o *Player, x float64, y float64) *RailNode {
	return &RailNode{
		Model:   NewModel(id),
		Owner:   NewOwner(o),
		Point:   NewPoint(x, y),
		InEdge:  make(map[uint]*RailEdge),
		OutEdge: make(map[uint]*RailEdge),
	}
}

// Idx returns unique id field.
func (rn *RailNode) Idx() uint {
	return rn.ID
}

// Pos returns location
func (rn *RailNode) Pos() *Point {
	return &rn.Point
}

// Out returns where it can go to
func (rn *RailNode) Out() map[uint]*Step {
	return nil
}

// In returns where it comes from
func (rn *RailNode) In() map[uint]*Step {
	return nil
}

// IsIn returns it should be view or not.
func (rn *RailNode) IsIn(center *Point, scale float64) bool {
	return rn.Pos().IsIn(center, scale)
}

// Resolve set reference
func (rn *RailNode) Resolve(owner *Player) {
	rn.Own = owner
	rn.ResolveRef()
}

// ResolveRef set id from reference
func (rn *RailNode) ResolveRef() {
	rn.Owner.ResolveRef()
}

// Permits represents Player is permitted to control
func (rn *RailNode) Permits(o *Player) bool {
	return rn.Owner.Permits(o)
}
