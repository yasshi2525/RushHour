package entities

// Indexable has unique id field.
type Indexable interface {
	// Idx returns unique id field.
	Idx() uint
}

// Locationable represents physical space.
type Locationable interface {
	// Pos returns entities' position.
	Pos() *Point
	// In represents how other can reach itself.
	In() map[uint]*Step
	// Out represents how itselt can reach other.
	Out() map[uint]*Step
	// IsIn returns it should be view or not.
	IsIn(*Point, float64) bool
}

// Connectable represents logical connection.
type Connectable interface {
	// From represents start point
	From() Locationable
	// To represents end point
	To() Locationable
	Cost() float64
}

// Resolvable set some_id fields from reference.
// Resolvable is for database migration
type Resolvable interface {
	// ResolveRef set id from object
	ResolveRef()
}

// Ownable is control of auth level
type Ownable interface {
	// Permits represents Player is permitted to control
	Permits(*Player) bool
}
