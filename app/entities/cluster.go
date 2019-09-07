package entities

import (
	"fmt"
	"math"
	"strings"

	"github.com/revel/revel"
)

const (
	// N represents North
	N = 0
	// S represents South
	S = 1
	// W represents West
	W = 0
	// E represents East
	E = 1
)

// Cluster summarize entities of each owner.
type Cluster struct {
	Base
	Point
	Scale  float64
	Parent *Cluster
	Data   map[uint]*Chunk

	ChPos    [2][2]*Point
	Children [2][2]*Cluster
}

// NewCluster creates Cluster on specified point.
func (m *Model) NewCluster(p *Cluster, dx int, dy int) *Cluster {
	cl := &Cluster{
		Base: m.NewBase(CLUSTER),
	}
	if p == nil {
		cl.Scale = Const.MaxScale
	} else {
		cl.Parent = p
		cl.Scale = p.Scale - 1
		len := math.Pow(2, p.Scale-2)
		cl.X = p.X + len*float64(dx)
		cl.Y = p.Y + len*float64(dy)
		x := int(math.Ceil(float64(dx) / 2))
		y := int(math.Ceil(float64(dy) / 2))
		p.Children[y][x] = cl
	}
	cl.Init(m)
	m.Add(cl)
	return cl
}

// B returns base information of this elements.
func (cl *Cluster) B() *Base {
	return &cl.Base
}

// Init creates map.
func (cl *Cluster) Init(m *Model) {
	cl.Base.Init(CLUSTER, m)

	cl.Data = make(map[uint]*Chunk)

	if cl.Scale > Const.MinScale {
		len := math.Pow(2, cl.Scale-2)
		cl.ChPos = [2][2]*Point{}
		cl.ChPos[N][W] = &Point{cl.X - len, cl.Y - len}
		cl.ChPos[N][E] = &Point{cl.X + len, cl.Y - len}
		cl.ChPos[S][W] = &Point{cl.X - len, cl.Y + len}
		cl.ChPos[S][E] = &Point{cl.X + len, cl.Y + len}
	}
}

// FindChunk returns specific scale's Chunk which has specified Entity.
func (cl *Cluster) FindChunk(obj Entity, scale float64) *Chunk {
	if cl.Scale == scale {
		data := cl.Data[obj.B().OwnerID]
		if data != nil && data.Has(obj) {
			return data
		}
	} else {
		for _, list := range cl.Children {
			for _, child := range list {
				if child != nil {
					if d := child.FindChunk(obj, scale); d != nil {
						return d
					}
				}
			}
		}
	}
	return nil
}

// FindChild returns child Cluster on specified point.
func (cl *Cluster) FindChild(dx int, dy int) (*Cluster, *Point) {
	x := int(math.Ceil(float64(dx) / 2))
	y := int(math.Ceil(float64(dy) / 2))
	return cl.Children[y][x], cl.ChPos[y][x]
}

// FindOrCreateChild must returns child Cluster on specified point.
func (cl *Cluster) FindOrCreateChild(dx int, dy int) *Cluster {
	if c, _ := cl.FindChild(dx, dy); c != nil {
		return c
	}
	return cl.M.NewCluster(cl, dx, dy)
}

// Add deploy Entity over related Chunk.
func (cl *Cluster) Add(raw Entity) {
	switch obj := raw.(type) {
	case *RailEdge:
		cl.addEntity(obj, obj.FromNode.Pos())
	}

	if obj, ok := raw.(Localable); ok {
		cl.addEntity(raw, obj.Pos())
	}
}

func (cl *Cluster) addEntity(obj Entity, p *Point) {
	if p == nil {
		return
	}

	if !cl.Point.IsIn(p.X, p.Y, cl.Scale) {
		revel.AppLog.Warnf("%v(%v) is out of bounds for %v", obj, p, cl)
	}

	oid := obj.B().OwnerID
	if _, ok := cl.Data[oid]; !ok {
		cl.Data[oid] = cl.M.NewChunk(cl, obj.B().O)
	}

	cl.Data[oid].Add(obj)

	cl.eachChildren(func(dx int, dy int, c *Cluster, pos *Point) {
		if pos.IsIn(p.X, p.Y, cl.Scale-1) {
			if c == nil {
				c = cl.M.NewCluster(cl, dx, dy)
			}
			c.Add(obj)
		}
	})
}

// Update changes Chunk of specified Entity
func (cl *Cluster) Update(obj Entity) {
	cl.Remove(obj)
	cl.Add(obj)
}

// Remove undeploy specified Entity over related Chunk.
func (cl *Cluster) Remove(raw Entity) {
	switch obj := raw.(type) {
	case *Cluster:
	case *Chunk:
	case *Player:
	case *RailLine:
	case *LineTask:
	case *Step:
	case *Track:
	case *Transport:
	default:
		oid := obj.B().OwnerID
		if chunk := cl.Data[oid]; chunk != nil {
			chunk.Remove(obj)
			cl.eachChildren(func(dx int, dy int, c *Cluster, p *Point) {
				if c != nil && c.Data[oid] != nil && c.Data[oid].Has(obj) {
					c.Data[oid].Remove(obj)
				}
			})
			if chunk.IsEmpty() {
				chunk.Delete()
			}
		}
		if len(cl.Data) == 0 {
			cl.Delete()
		}
	}
}

// ViewMap set delegate Entity to DelegateMap.
func (cl *Cluster) ViewMap(dm *DelegateMap, cx float64, cy float64, scale float64, span float64) {
	if cl.intersectsWith(cx, cy, scale) {
		if cl.Scale <= scale-span {
			for _, d := range cl.Data {
				d.Export(dm)
			}
		} else {
			cl.eachChildren(func(dx int, dy int, c *Cluster, p *Point) {
				if c != nil {
					c.ViewMap(dm, cx, cy, scale, span)
				}
			})
		}
	}
}

func (cl *Cluster) eachChildren(callback func(int, int, *Cluster, *Point)) {
	if cl.Scale > Const.MinScale {
		for _, dy := range []int{-1, +1} {
			for _, dx := range []int{-1, +1} {
				c, p := cl.FindChild(dx, dy)
				callback(dx, dy, c, p)
			}
		}
	}
}

// BeforeDelete remove reference of related entity
func (cl *Cluster) BeforeDelete() {
	if cl.Parent != nil {
		cl.Parent.UnResolve(cl)
	}
}

// UnResolve unregisters specified refernce.
func (cl *Cluster) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Cluster:
			for y, list := range cl.Children {
				for x, child := range list {
					if child == cl {
						cl.Children[y][x] = nil
					}
				}
			}
		case *Chunk:
			delete(cl.Data, obj.OwnerID)
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Resolve set reference.
func (cl *Cluster) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Chunk:
			cl.Data[obj.OwnerID] = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// CheckDelete check remain relation.
func (cl *Cluster) CheckDelete() error {
	if len(cl.Data) > 0 {
		return fmt.Errorf("data exists")
	}
	if cl.Scale > Const.MinScale {
		for _, dy := range []int{-1, +1} {
			for _, dx := range []int{-1, +1} {
				if res := cl.Children[dy][dx].CheckDelete(); res != nil {
					return res
				}
			}
		}
	}
	return nil
}

// Delete removes this entity with related ones.
func (cl *Cluster) Delete() {
	for _, list := range cl.Children {
		for _, child := range list {
			if child != nil {
				child.Delete()
			}
		}
	}
	for _, ch := range cl.Data {
		ch.Delete()
	}
	cl.M.Delete(cl)
}

func (cl *Cluster) intersectsWith(cx float64, cy float64, scale float64) bool {
	myL := math.Pow(2, cl.Scale) / 2
	othL := math.Pow(2, scale) / 2

	return math.Max(cl.X-myL, cx-othL) <= math.Min(cl.X+myL, cx+othL) &&
		math.Max(cl.Y-myL, cy-othL) <= math.Min(cl.Y+myL, cy+othL)
}

// String represents status
func (cl *Cluster) String() string {
	list := []string{}
	for id := range cl.Data {
		list = append(list, fmt.Sprintf("ch(%d)", id))
	}
	return fmt.Sprintf("%s(%.1f:%d):%s,%v", cl.Type().Short(),
		cl.Scale, cl.ID, strings.Join(list, ","), cl.Point)
}
