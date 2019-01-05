package entities

import "time"

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
	Idx() uint
	Type() ModelType
	// Pos returns entities' position.
	Pos() *Point
}

// Relayable represents connectable for Human moving
type Relayable interface {
	Idx() uint
	Type() ModelType
	// Pos returns entities' position.
	Pos() *Point
	// In represents how other can reach itself.
	InStep() map[uint]*Step
	// Out represents how itselt can reach other.
	OutStep() map[uint]*Step
}

// Connectable represents movability of two resource
type Connectable interface {
	Idx() uint
	Type() ModelType
	From() Indexable
	To() Indexable
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

// Persistable represents Differencial backup
type Persistable interface {
	IsChanged(after ...time.Time) bool
	Reset()
}

// Viewable represents user can view it
type Viewable interface {
	// Pos returns entities' position.
	Pos() *Point
	// IsIn returns it should be view or not.
	IsIn(float64, float64, float64) bool
}

// Removable can check to able to remove
type Removable interface {
	CheckRemove() error
}
