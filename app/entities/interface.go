package entities


// Entity represents that it is CRUD object.
type Entity interface {
	B() *Base
	S() *Shape
	CheckDelete() error
	BeforeDelete()
	Delete(bool)
	Resolve(...Entity)
}

// Initializable represents that setup is required.
type Initializable interface {
	// Init will be called after instanciation.
	Init(*Model)
}

// Relayable represents connectable for Human moving
type Relayable interface {
	// In represents how other can reach itself.
	InSteps() map[uint]*Step
	// Out represents how itselt can reach other.
	OutSteps() map[uint]*Step

	B() *Base
	S() *Shape
	CheckDelete() error
	BeforeDelete()
	Delete(bool)
	Resolve(...Entity)
}

// Connectable represents movability of two resource
type Connectable interface {
	B() *Base
	From() Entity
	To() Entity
	Cost() float64
}

// Migratable set some_id fields from reference.
// Migratable is for database migration
type Migratable interface {
	// Marshal set id from object
	Marshal()
	// Unmarshal set object from id
	UnMarshal()
}

// Persistable represents Differencial backup
type Persistable interface {
	B() *Base
	P() *Persistence
}

// Steppable represents human over it can step forward at certain rate.
type Steppable interface {
	Step(*float64, *float64)
	Loc(float64) *Point
}
