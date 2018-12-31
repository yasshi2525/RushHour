package entities

import "github.com/jinzhu/gorm"

// Resolvable set some_id fields from reference.
// Resolvable is for database migration
type Resolvable interface {
	ResolveRef()
}

// Ownable means this faciliites in under the control by Player.
type Ownable struct {
	Owner *Player

	OwnerID uint `gorm:"not null"`
}

// NewOwnable create Juntion
func NewOwnable(o *Player) Ownable {
	return Ownable{
		Owner:   o,
		OwnerID: o.ID,
	}
}

// ResolveRef resolve ownerID from Owner
func (o *Ownable) ResolveRef() {
	o.OwnerID = o.Owner.ID
}

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	gorm.Model
	Ownable
	Point

	In  []*RailEdge `gorm:"-"`
	Out []*RailEdge `gorm:"-"`
}

// ResolveRef resolve Owner reference
func (rn *RailNode) ResolveRef() {
	rn.Ownable.ResolveRef()
}

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	gorm.Model
	Ownable

	From *RailNode `gorm:"-"`
	To   *RailNode `gorm:"-"`

	FromID uint `gorm:"not null"`
	ToID   uint `gorm:"not null"`
}

// ResolveRef resolve Owner reference
func (re *RailEdge) ResolveRef() {
	re.Ownable.ResolveRef()
	re.FromID = re.From.ID
	re.ToID = re.To.ID
}

// Station composes on Platform and Gate
type Station struct {
	gorm.Model
	Ownable

	Name     string
	Platform *Platform `gorm:"-"`
	Gate     *Gate     `gorm:"-"`
}

// ResolveRef resolve Owner reference
func (st *Station) ResolveRef() {
	st.Ownable.ResolveRef()
}

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	gorm.Model
	Ownable
	Junction

	In       *Station  `gorm:"-"`
	On       *RailNode `gorm:"-"`
	Capacity uint      `gorm:"not null"`
	Occupied uint      `gorm:"not null"`

	InID uint `gorm:"not null"`
	OnID uint `gorm:"not null"`
}

// ResolveRef resolve Owner and Station reference
func (p *Platform) ResolveRef() {
	p.Ownable.ResolveRef()
	p.InID = p.In.ID
	p.OnID = p.On.ID
}

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	gorm.Model
	Ownable
	Junction

	In *Station `gorm:"-"`
	// Num represents how many Human can pass at the same time
	Num uint `gorm:"not null"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null"`
	// Occupied represents how many Gate are used by Human.
	Occupied uint `gorm:"not null"`

	InID uint `gorm:"not null"`
}

// ResolveRef resolve Owner and Station reference
func (g *Gate) ResolveRef() {
	g.Ownable.ResolveRef()
	g.InID = g.In.ID
}

// LineTaskType represents the state what Train should do now.
type LineTaskType uint

const (
	// OnDeparture represents the state that Train waits for departure in Station.
	OnDeparture LineTaskType = iota
	// OnMoving represents the state that Train runs to next RailNode.
	OnMoving
	// OnStopping represents the state that Train stops to next Station.
	OnStopping
	// OnPassing represents the state that Train passes to next Station.
	OnPassing
)

// LineTask is the element of Line.
// The chain of LineTask represents Line structure.
type LineTask struct {
	gorm.Model
	Ownable

	Line *Line        `gorm:"-"`
	Type LineTaskType `gorm:"not null"`
	Next *LineTask    `gorm:"-"`

	LineID uint `gorm:"not null"`
	NextID uint
}

// ResolveRef resolve Owner and Line reference
func (lt *LineTask) ResolveRef() {
	lt.Ownable.ResolveRef()
	lt.LineID = lt.Line.ID
	lt.NextID = lt.Next.ID
}

// Line represents how Train should run.
type Line struct {
	gorm.Model
	Ownable

	Name  string
	Tasks []*LineTask `gorm:"-"`
}

// ResolveRef resolve Owner reference
func (l *Line) ResolveRef() {
	l.Ownable.ResolveRef()
}

// Train carries Human from Station to Station.
type Train struct {
	gorm.Model
	Ownable
	Point

	Capacity uint `gorm:"not null"`
	// Mobility represents how many Human can get off at the same time.
	Mobility uint      `gorm:"not null"`
	Speed    float64   `gorm:"not null"`
	Name     string    `gorm:"not null"`
	Task     *LineTask `gorm:"-"`

	TaskID uint `gorm:"not null"`
}

// ResolveRef resolve Owner reference
func (t *Train) ResolveRef() {
	t.Ownable.ResolveRef()
	t.TaskID = t.Task.ID
}
