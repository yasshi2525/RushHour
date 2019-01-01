package entities

// Train carries Human from Station to Station.
type Train struct {
	Model
	Owner

	Capacity uint `gorm:"not null"`
	// Mobility represents how many Human can get off at the same time.
	Mobility uint    `gorm:"not null"`
	Speed    float64 `gorm:"not null"`
	Name     string  `gorm:"not null"`
	Progress float64 `gorm:"not null"`

	Task      *LineTask       `gorm:"-" json:"-"`
	Passenger map[uint]*Human `gorm:"-" json:"-"`

	TaskID uint
}

// NewTrain creates instance
func NewTrain(id uint, o *Player) *Train {
	return &Train{
		Model:     NewModel(id),
		Owner:     NewOwner(o),
		Passenger: make(map[uint]*Human),
	}
}

// Idx returns unique id field.
func (t *Train) Idx() uint {
	return t.ID
}

// Pos returns location
func (t *Train) Pos() *Point {
	if t.Task == nil {
		return nil
	}
	from, to := t.Task.From().Pos(), t.Task.To().Pos()
	return from.Div(to, t.Progress)
}

// Out returns where it can go to
func (t *Train) Out() map[uint]*Step {
	return nil
}

// In returns where it comes from
func (t *Train) In() map[uint]*Step {
	return nil
}

// IsIn returns it should be view or not.
func (t *Train) IsIn(center *Point, scale float64) bool {
	return t.Pos().IsIn(center, scale)
}

// Resolve set reference
func (t *Train) Resolve(lt *LineTask) {
	t.Owner, t.Task = lt.Owner, lt
	t.ResolveRef()
}

// ResolveRef set id from reference
func (t *Train) ResolveRef() {
	t.Owner.ResolveRef()
	t.TaskID = t.Task.ID
}

// Permits represents Player is permitted to control
func (t *Train) Permits(o *Player) bool {
	return t.Owner.Permits(o)
}
