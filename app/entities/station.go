package entities

// Station composes on Platform and Gate
type Station struct {
	Model
	Owner

	Platform *Platform `gorm:"-" json:"-"`
	Gate     *Gate     `gorm:"-" json:"-"`

	Name string
}

// NewStation create new instance.
func NewStation(stid uint, gid uint, pid uint, rn *RailNode) (*Station, *Gate, *Platform) {
	p := &Platform{
		Model:      NewModel(pid),
		Owner:      rn.Owner,
		in:         make(map[uint]*Step),
		out:        make(map[uint]*Step),
		OnRailNode: rn,
	}
	rn.OverPlatform = p

	g := &Gate{
		Model: NewModel(gid),
		Owner: rn.Owner,
		in:    make(map[uint]*Step),
		out:   make(map[uint]*Step),
	}

	st := &Station{
		Model:    NewModel(stid),
		Owner:    rn.Owner,
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

// Idx returns unique id field.
func (st *Station) Idx() uint {
	return st.ID
}

// Init creates map.
func (st *Station) Init() {
	st.Model.Init()
	st.Owner.Init()
}

// Pos returns location
func (st *Station) Pos() *Point {
	return st.Platform.Pos()
}

// IsIn returns it should be view or not.
func (st *Station) IsIn(center *Point, scale float64) bool {
	return st.Pos().IsIn(center, scale)
}

// Resolve set reference from id.
func (st *Station) Resolve(o *Player) {
	st.Own = o
	st.ResolveRef()
}

// ResolveRef resolve Owner reference
func (st *Station) ResolveRef() {
	st.Owner.ResolveRef()
}

// Permits represents Player is permitted to control
func (st *Station) Permits(o *Player) bool {
	return st.Owner.Permits(o)
}
