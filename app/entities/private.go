package entities

// RailNode represents rail track as point.
// Station only stands on RailNode.
type RailNode struct {
	Model
	Ownable

	Loc *Point
	In  map[uint]*RailEdge `gorm:"-" json:"-"`
	Out map[uint]*RailEdge `gorm:"-" json:"-"`

	OverPlatform *Platform `gorm:"-" json:"-"`
}

// NewRailNode create new instance.
func NewRailNode(id uint, owner *Player, x float64, y float64) *RailNode {
	return &RailNode{
		Model:   NewModel(id),
		Ownable: NewOwnable(owner),
		Loc:     NewPoint(x, y),
		In:      make(map[uint]*RailEdge),
		Out:     make(map[uint]*RailEdge),
	}
}

// Resolve set reference
func (rn *RailNode) Resolve(owner *Player) {
	rn.Owner = owner
	rn.ResolveRef()
}

// ResolveRef set id from reference
func (rn *RailNode) ResolveRef() {
	rn.Ownable.ResolveRef()
}

// Permits represents Player is permitted to control
func (rn *RailNode) Permits(o *Player) bool {
	return rn.Ownable.Permits(o)
}

// RailEdge connects from RailNode to RailNode.
// It's directional.
type RailEdge struct {
	Model
	Ownable
	Loc  *Point
	From *RailNode `gorm:"-" json:"-"`
	To   *RailNode `gorm:"-" json:"-"`

	FromID uint `gorm:"not null"`
	ToID   uint `gorm:"not null"`
}

// NewRailEdge create new instance and relates RailNode
func NewRailEdge(id uint, from *RailNode, to *RailNode) *RailEdge {
	re := &RailEdge{
		Model:   NewModel(id),
		Ownable: from.Ownable,
		Loc:     from.Loc.Center(to.Loc),
		From:    from,
		To:      to,
	}
	re.ResolveRef()

	from.Out[re.ID] = re
	to.In[re.ID] = re
	return re
}

// Unrelate delete relations to RailNode
func (re *RailEdge) Unrelate() {
	delete(re.From.Out, re.ID)
	delete(re.To.In, re.ID)
}

// Resolve set reference.
func (re *RailEdge) Resolve(from *RailNode, to *RailNode) {
	re.Owner, re.From, re.To = from.Owner, from, to

	from.Out[re.ID] = re
	to.In[re.ID] = re

	re.ResolveRef()
}

// ResolveRef set id from reference
func (re *RailEdge) ResolveRef() {
	re.Ownable.ResolveRef()
	re.FromID = re.From.ID
	re.ToID = re.To.ID
}

// Permits represents Player is permitted to control
func (re *RailEdge) Permits(o *Player) bool {
	return re.Ownable.Permits(o)
}

// Station composes on Platform and Gate
type Station struct {
	Model
	Ownable
	Loc *Point

	Platform *Platform `gorm:"-" json:"-"`
	Gate     *Gate     `gorm:"-" json:"-"`

	Name string
}

// NewStation create new instance.
func NewStation(stid uint, gid uint, pid uint, rn *RailNode) (*Station, *Gate, *Platform) {
	p := &Platform{
		Model:   NewModel(pid),
		Ownable: rn.Ownable,

		Loc: rn.Loc,
		In:  make(map[uint]*Step),
		Out: make(map[uint]*Step),

		OnRailNode: rn,
	}
	rn.OverPlatform = p

	g := &Gate{
		Model:   NewModel(gid),
		Ownable: rn.Ownable,

		Loc: rn.Loc,
		In:  make(map[uint]*Step),
		Out: make(map[uint]*Step),
	}

	st := &Station{
		Model:    NewModel(stid),
		Ownable:  rn.Ownable,
		Loc:      rn.Loc,
		Platform: p,
		Gate:     g,
	}

	p.InStation = st
	g.InStation = st

	p.ResolveRef()
	g.ResolveRef()
	st.ResolveRef()

	return st, g, p
}

// Resolve set reference from id.
func (st *Station) Resolve(o *Player) {
	st.Owner = o
	st.ResolveRef()
}

// ResolveRef resolve Owner reference
func (st *Station) ResolveRef() {
	st.Ownable.ResolveRef()
}

// Permits represents Player is permitted to control
func (st *Station) Permits(o *Player) bool {
	return st.Ownable.Permits(o)
}

// Platform is the base Human wait for Train.
// Platform can enter only through Gate.
type Platform struct {
	Model
	Ownable
	Loc *Point
	Out map[uint]*Step `gorm:"-" json:"-"`
	In  map[uint]*Step `gorm:"-" json:"-"`

	InStation  *Station  `gorm:"-" json:"-"`
	OnRailNode *RailNode `gorm:"-" json:"-"`
	Passenger  []*Human  `gorm:"-" json:"-"`

	Capacity uint `gorm:"not null"`
	Occupied uint `gorm:"not null"`

	InStationID  uint `gorm:"not null"`
	OnRailNodeID uint `gorm:"not null"`
}

// Resolve set reference
func (p *Platform) Resolve(rn *RailNode, st *Station) {
	p.Owner, p.OnRailNode, p.InStation = rn.Owner, rn, st
	p.ResolveRef()
}

// ResolveRef set id from reference
func (p *Platform) ResolveRef() {
	p.Ownable.ResolveRef()
	p.OnRailNodeID = p.OnRailNode.ID
	p.InStationID = p.InStation.ID
}

// Pos returns location
func (p *Platform) Pos() *Point {
	return p.Loc
}

// OutStep returns where it can go to
func (p *Platform) OutStep() map[uint]*Step {
	return p.Out
}

// InStep returns where it comes from
func (p *Platform) InStep() map[uint]*Step {
	return p.In
}

// Permits represents Player is permitted to control
func (p *Platform) Permits(o *Player) bool {
	return p.Ownable.Permits(o)
}

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	Model
	Ownable

	Loc *Point
	Out map[uint]*Step `gorm:"-" json:"-"`
	In  map[uint]*Step `gorm:"-" json:"-"`

	InStation *Station `gorm:"-" json:"-"`

	// Num represents how many Human can pass at the same time
	Num uint `gorm:"not null"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null"`
	// Occupied represents how many Gate are used by Human.
	Occupied uint `gorm:"not null"`

	InStationID uint `gorm:"not null"`
}

// Resolve set reference
func (g *Gate) Resolve(st *Station) {
	g.Owner, g.InStation = st.Owner, st
	g.ResolveRef()
}

// ResolveRef set id from reference
func (g *Gate) ResolveRef() {
	g.Ownable.ResolveRef()
	g.InStationID = g.InStation.ID
}

// Pos returns location
func (g *Gate) Pos() *Point {
	return g.Loc
}

// OutStep returns where it can go to
func (g *Gate) OutStep() map[uint]*Step {
	return g.Out
}

// InStep returns where it comes from
func (g *Gate) InStep() map[uint]*Step {
	return g.In
}

// SetOwner set owner object and owner id
func (g *Gate) SetOwner(o *Player) {
	g.Owner = o
	g.ResolveRef()
}

// Permits represents Player is permitted to control
func (g *Gate) Permits(o *Player) bool {
	return g.Ownable.Permits(o)
}

// Line represents how Train should run.
type Line struct {
	Model
	Ownable

	Name  string
	Tasks []*LineTask `gorm:"-" json:"-"`
}

// Resolve set reference
func (l *Line) Resolve(o *Player) {
	l.Owner = o
	l.ResolveRef()
}

// ResolveRef set if from reference
func (l *Line) ResolveRef() {
	l.Ownable.ResolveRef()
}

// Permits represents Player is permitted to control
func (l *Line) Permits(o *Player) bool {
	return l.Ownable.Permits(o)
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
	Model
	Ownable

	Line *Line        `gorm:"-" json:"-"`
	Type LineTaskType `gorm:"not null"`
	Next *LineTask    `gorm:"-" json:"-"`

	LineID uint `gorm:"not null"`
	NextID uint
}

// Resolve set reference
func (lt *LineTask) Resolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Line:
			lt.Owner, lt.Line = obj.Owner, obj
		case *LineTask:
			lt.Next = obj
		}
	}

	lt.ResolveRef()
}

// ResolveRef set id from reference
func (lt *LineTask) ResolveRef() {
	lt.Ownable.ResolveRef()
	lt.LineID = lt.Line.ID
	if lt.Next != nil {
		lt.NextID = lt.Next.ID
	}
}

func (lt *LineTask) Permits(o *Player) bool {
	return lt.Ownable.Permits(o)
}

// Train carries Human from Station to Station.
type Train struct {
	Model
	Ownable
	Point

	Capacity uint `gorm:"not null"`
	// Mobility represents how many Human can get off at the same time.
	Mobility uint    `gorm:"not null"`
	Speed    float64 `gorm:"not null"`
	Name     string  `gorm:"not null"`
	Progress float64 `gorm:"not null"`

	Task      *LineTask `gorm:"-" json:"-"`
	Passenger []*Human  `gorm:"-" json:"-"`

	TaskID uint `gorm:"not null"`
}

// Resolve set reference
func (t *Train) Resolve(lt *LineTask) {
	t.Owner, t.Task = lt.Owner, lt
	t.ResolveRef()
}

// ResolveRef set id from reference
func (t *Train) ResolveRef() {
	t.Ownable.ResolveRef()
	t.TaskID = t.Task.ID
}

// Permits represents Player is permitted to control
func (t *Train) Permits(o *Player) bool {
	return t.Ownable.Permits(o)
}
