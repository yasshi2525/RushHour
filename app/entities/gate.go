package entities

import (
	"fmt"
)

// Gate represents ticket gate in Station.
// Human must pass Gate to enter/leave Platform.
type Gate struct {
	Base
	Persistence
	Shape

	// Num represents how many Human can pass at the same time
	Num int `gorm:"not null" json:"num"`
	// Mobility represents time one Human pass Gate.
	Mobility float64 `gorm:"not null" json:"mobility"`
	// Occupied represents how many Gate are used by Human.
	Occupied int `gorm:"not null" json:"occupied"`

	InStation    *Station  `gorm:"-" json:"-"`
	WithPlatform *Platform `gorm:"-" json:"-"`
	out          map[uint]*Step
	in           map[uint]*Step

	StationID  uint `gorm:"not null" json:"stid"`
	PlatformID uint `gorm:"-"        json:"pid"`
}

// NewGate creates instance
func (m *Model) NewGate(st *Station) *Gate {
	g := &Gate{
		Base:        m.NewBase(GATE, st.O),
		Persistence: NewPersistence(),
		Num:         Const.Gate.Num,
	}
	g.Init(m)
	g.Resolve(st.O, st)
	g.Marshal()
	m.Add(g)

	g.GenOutSteps()
	g.GenInSteps()

	return g
}

// B returns base information of this elements.
func (g *Gate) B() *Base {
	return &g.Base
}

// P returns time information for database.
func (g *Gate) P() *Persistence {
	return &g.Persistence
}

// S returns entities' position.
func (g *Gate) S() *Shape {
	return &g.Shape
}

// GenOutSteps generates Steps from this Gate.
func (g *Gate) GenOutSteps() {
	// skip G -> P
	// G -> C
	for _, c := range g.M.Companies {
		g.M.NewStep(g, c)
	}
}

// GenInSteps generates Steps to this Gate.
func (g *Gate) GenInSteps() {
	// skip P -> G
	// R -> G
	for _, r := range g.M.Residences {
		g.M.NewStep(r, g)
	}
}

// Init creates map.
func (g *Gate) Init(m *Model) {
	g.Base.Init(GATE, m)
	g.out = make(map[uint]*Step)
	g.in = make(map[uint]*Step)
}

// Pos returns location
func (g *Gate) Pos() *Point {
	if g.WithPlatform == nil {
		return nil
	}
	return g.WithPlatform.Pos()
}

// IsIn returns it should be view or not.
func (g *Gate) IsIn(x float64, y float64, scale float64) bool {
	return g.Pos().IsIn(x, y, scale)
}

// OutSteps returns where it can go to
func (g *Gate) OutSteps() map[uint]*Step {
	return g.out
}

// InSteps returns where it comes from
func (g *Gate) InSteps() map[uint]*Step {
	return g.in
}

// Resolve set reference
func (g *Gate) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Player:
			g.O = obj
			obj.Resolve(g)
		case *Station:
			g.InStation = obj
			obj.Resolve(g)
		case *Platform:
			g.WithPlatform = obj
			g.Shape = obj.Shape
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
	g.Marshal()
}

// Marshal set id from reference
func (g *Gate) Marshal() {
	if g.InStation != nil {
		g.StationID = g.InStation.ID
	}
	if g.WithPlatform != nil {
		g.PlatformID = g.WithPlatform.ID
	}
}

// UnMarshal set reference from id.
func (g *Gate) UnMarshal() {
	g.Resolve(
		g.M.Find(PLAYER, g.OwnerID),
		g.M.Find(STATION, g.StationID))
}

// CheckDelete check remain relation.
func (g *Gate) CheckDelete() error {
	return nil
}

// BeforeDelete deletes related reference
func (g *Gate) BeforeDelete() {
	g.O.UnResolve(g)
}

// Delete removes this entity with related ones.
func (g *Gate) Delete(force bool) {
	for _, s := range g.out {
		g.M.Delete(s)
	}
	for _, s := range g.in {
		g.M.Delete(s)
	}
	g.M.Delete(g)
}

// String represents status
func (g *Gate) String() string {
	g.Marshal()
	ostr := ""
	if g.O != nil {
		ostr = fmt.Sprintf(":%s", g.O.Short())
	}
	ststr := ""
	if g.InStation != nil {
		ststr = fmt.Sprintf(":%s", g.InStation.Name)
	}
	posstr := ""
	if g.Pos() != nil {
		posstr = fmt.Sprintf(":%s", g.Pos())
	}
	return fmt.Sprintf("%s(%d):st=%d,p=%d,i=%d,o=%d%s%s%s", g.Type().Short(),
		g.ID, g.StationID, g.PlatformID, len(g.in), len(g.out), posstr, ostr, ststr)
}
