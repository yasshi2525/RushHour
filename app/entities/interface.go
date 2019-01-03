package entities

// Indexable has unique id field.
type Indexable interface {
	// Idx returns unique id field.
	Idx() uint
	Type() ModelType
}

// Initializable represents that setup is required.
type Initializable interface {
	// Init will be called after instanciation.
	Init()
}

// Locationable represents physical space.
type Locationable interface {
	// Pos returns entities' position.
	Pos() *Point
	// IsIn returns it should be view or not.
	IsIn(*Point, float64) bool
}

// Relayable represents connectable for Human moving
type Relayable interface {
	// Pos returns entities' position.
	Pos() *Point
	// In represents how other can reach itself.
	In() map[uint]*Step
	// Out represents how itselt can reach other.
	Out() map[uint]*Step
}

// Connectable represents logical connection.
type Connectable interface {
	// From represents start point
	From() Relayable
	// To represents end point
	To() Relayable
	Cost() float64
}

// Resolvable set some_id fields from reference.
// Resolvable is for database migration
type Resolvable interface {
	// ResolveRef set id from object
	ResolveRef()
}

// UnReferable can unrefer for related entity.
type UnReferable interface {
	Idx() uint
	Type() ModelType
	UnRef()
}

// Ownable is control of auth level
type Ownable interface {
	// Permits represents Player is permitted to control
	Permits(*Player) bool
}
