package entities

// RailLine represents how Train should run.
type RailLine struct {
	Model
	Owner

	Name  string
	Tasks map[uint]*LineTask `gorm:"-" json:"-"`
}

// NewRailLine create instance
func NewRailLine(id uint, o *Player) *RailLine {
	return &RailLine{
		Model: NewModel(id),
		Owner: NewOwner(o),
		Tasks: make(map[uint]*LineTask),
	}
}

// Idx returns unique id field.
func (l *RailLine) Idx() uint {
	return l.ID
}

// Init makes map
func (l *RailLine) Init() {
	l.Model.Init()
	l.Owner.Init()
	l.Tasks = make(map[uint]*LineTask)
}

// Pos returns location
func (l *RailLine) Pos() *Point {
	sumX, sumY, cnt := 0.0, 0.0, 0.0
	for _, lt := range l.Tasks {
		sumX += lt.Pos().X
		sumY += lt.Pos().Y
		cnt++
	}
	if cnt > 0 {
		return &Point{sumX / cnt, sumY / cnt}
	}
	return nil
}

// IsIn return true when any LineTask is in,
func (l *RailLine) IsIn(center *Point, scale float64) bool {
	for _, lt := range l.Tasks {
		if lt.IsIn(center, scale) {
			return true
		}
	}
	return false
}

// Resolve set reference
func (l *RailLine) Resolve(o *Player) {
	l.Own = o
	l.ResolveRef()
}

// ResolveRef set if from reference
func (l *RailLine) ResolveRef() {
	l.Owner.ResolveRef()
}

// Permits represents Player is permitted to control
func (l *RailLine) Permits(o *Player) bool {
	return l.Owner.Permits(o)
}
