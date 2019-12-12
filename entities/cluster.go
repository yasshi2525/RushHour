package entities

import (
	"fmt"
	"log"
	"strings"
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

// ChunkPoint is scaled value
type ChunkPoint struct {
	X     int
	Y     int
	Scale int
}

func (p *ChunkPoint) has(pos *Point) bool {
	x, y := int(pos.X), int(pos.Y)
	left, right := p.X<<p.Scale, p.X<<p.Scale+1
	top, buttom := p.Y<<p.Scale, p.Y<<p.Scale+1
	return left <= x && x < right && top <= y && y < buttom
}

func (p *ChunkPoint) contains(oth *ChunkPoint) bool {
	if p.Scale < oth.Scale {
		return false
	}
	diff := p.Scale - oth.Scale
	return p.X == oth.X>>diff && p.Y == oth.Y>>diff
}

func (p *ChunkPoint) String() string {
	return fmt.Sprintf("(%d,%d,%d)", p.X, p.Y, p.Scale)
}

// Cluster summarize entities of each owner.
type Cluster struct {
	Base
	ChunkPoint
	Parent *Cluster
	Data   map[uint]*Chunk

	Children [2][2]*Cluster
}

// NewCluster creates Cluster on specified point.
// dx, dy must be 0 or 1
func (m *Model) NewCluster(p *Cluster, dx int, dy int) *Cluster {
	cl := &Cluster{
		Base: m.NewBase(CLUSTER),
	}
	if p == nil {
		cl.Scale = m.conf.MaxScale
	} else {
		cl.Parent = p
		cl.Scale = p.Scale - 1
		cl.X = p.X<<1 + dx
		cl.Y = p.Y<<1 + dy
		p.Children[dy][dx] = cl
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
}

// FindChunk returns specific scale's Chunk which has specified Entity.
func (cl *Cluster) FindChunk(obj Entity, scale int) *Chunk {
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

// Add deploy Entity over related Chunk.
func (cl *Cluster) Add(raw Entity) {
	if obj, ok := raw.(Connectable); ok {
		cl.addEntity(raw, obj.From().(Localable).Pos())
	}
	if obj, ok := raw.(Localable); ok {
		cl.addEntity(raw, obj.Pos())
	}
}

func (cl *Cluster) addEntity(obj Entity, p *Point) {
	if p == nil {
		return
	}
	if cl.Parent == nil && (int(p.X)>>cl.Scale > 0 || int(p.Y)>>cl.Scale > 0) {
		log.Printf("%v(%v) is out of bounds for %v", obj, p, cl)
	}

	oid := obj.B().OwnerID
	if _, ok := cl.Data[oid]; !ok {
		cl.Data[oid] = cl.M.NewChunk(cl, obj.B().O)
	}

	cl.Data[oid].Add(obj)

	cl.eachChildren(func(dx int, dy int, c *Cluster, pos *ChunkPoint) {
		if pos.has(p) {
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
func (cl *Cluster) Remove(obj Entity) {
	if _, ok := obj.(Localable); ok {
		cl.removeEntity(obj)
	}
	if _, ok := obj.(Connectable); ok {
		cl.removeEntity(obj)
	}
}

func (cl *Cluster) removeEntity(obj Entity) {
	oid := obj.B().OwnerID
	if chunk := cl.Data[oid]; chunk != nil {
		chunk.Remove(obj)
		cl.eachChildren(func(dx int, dy int, c *Cluster, p *ChunkPoint) {
			if c != nil && c.Data[oid] != nil && c.Data[oid].Has(obj) {
				c.Remove(obj)
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

// ViewMap set delegate Entity to DelegateMap.
func (cl *Cluster) ViewMap(dm *DelegateMap, pos *ChunkPoint, span int) {
	if cl.ChunkPoint.contains(pos) {
		if cl.Scale <= pos.Scale-span {
			for _, d := range cl.Data {
				d.Export(dm)
			}
		} else {
			cl.eachChildren(func(dx int, dy int, c *Cluster, p *ChunkPoint) {
				if c != nil {
					c.ViewMap(dm, pos, span)
				}
			})
		}
	}
}

func (cl *Cluster) eachChildren(callback func(int, int, *Cluster, *ChunkPoint)) {
	if cl.Scale > cl.M.conf.MinScale {
		for _, dy := range []int{0, 1} {
			for _, dx := range []int{0, 1} {
				callback(dx, dy, cl.Children[dy][dx],
					&ChunkPoint{cl.X<<1 + dx, cl.Y<<1 + dy, cl.Scale - 1})
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
	if cl.Scale > cl.M.conf.MinScale {
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

// String represents status
func (cl *Cluster) String() string {
	list := []string{}
	for id := range cl.Data {
		list = append(list, fmt.Sprintf("ch(%d)", id))
	}
	return fmt.Sprintf("%s(%d:%d):%s,%v", cl.Type().Short(),
		cl.Scale, cl.ID, strings.Join(list, ","), cl.ChunkPoint)
}
