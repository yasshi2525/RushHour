package entities

import (
	"fmt"
	"math"
)

const (
	N = 0
	S = 1
	W = 0
	E = 1
)

// Cluster summarize entities of each owner.
type Cluster struct {
	Base
	Shape
	Point
	Scale  float64
	Parent *Cluster
	Data   map[uint]*Chunk

	ChPos    [2][2]*Point
	Children [2][2]*Cluster
}

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

func (cl *Cluster) B() *Base {
	return &cl.Base
}

func (cl *Cluster) S() *Shape {
	return &cl.Shape
}

func (cl *Cluster) Init(m *Model) {
	cl.Base.Init(CLUSTER, m)

	cl.Shape.P1 = &cl.Point

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

func (cl *Cluster) FindChunk(obj Entity, scale float64) *Chunk {
	if cl.Scale == scale {
		if data := cl.Data[obj.B().OwnerID]; data.Has(obj) {
			return data
		}
	} else {
		for _, list := range cl.Children {
			for _, child := range list {
				if d := child.FindChunk(obj, scale); d != nil {
					return d
				}
			}
		}
	}
	return nil
}

func (cl *Cluster) FindChild(dx int, dy int) *Cluster {
	x := int(math.Ceil(float64(dx) / 2))
	y := int(math.Ceil(float64(dy) / 2))
	return cl.Children[y][x]
}

func (cl *Cluster) FindOrCreateChild(dx int, dy int) *Cluster {
	if c := cl.FindChild(dx, dy); c != nil {
		return c
	}
	return cl.M.NewCluster(cl, dx, dy)
}

func (cl *Cluster) Add(raw Entity) {
	switch obj := raw.(type) {
	case *Cluster:
	case *Chunk:
	default:
		oid := obj.B().OwnerID
		if _, ok := cl.Data[oid]; !ok {
			cl.Data[oid] = cl.M.NewChunk(cl, obj.B().O)
		}

		cl.Data[oid].Add(obj)

		len := math.Pow(2, cl.Scale-2)
		cl.EachChildren(func(dx int, dy int, c *Cluster) {
			if obj.S().IsIn(cl.X+len*float64(dx), cl.Y+len*float64(dy), cl.Scale-1) {
				if c == nil {
					c = cl.M.NewCluster(cl, dx, dy)
				}
				c.Add(obj)
			}
		})
	}
}

func (cl *Cluster) Update(obj Entity) {
	cl.Remove(obj)
	cl.Add(obj)
}

func (cl *Cluster) Remove(raw Entity) {
	switch obj := raw.(type) {
	case *Cluster:
	case *Chunk:
	default:
		oid := obj.B().OwnerID
		if chunk := cl.Data[oid]; chunk != nil {
			chunk.Remove(obj)
			cl.EachChildren(func(dx int, dy int, c *Cluster) {
				if c != nil && c.Data[oid] != nil && c.Data[oid].Has(obj) {
					chunk.Remove(obj)
				}
			})
			if chunk.IsEmpty() {
				cl.Delete()
			}
		}
	}
}

func (cl *Cluster) ViewMap(dm *DelegateMap, cx float64, cy float64, scale float64, span float64) {
	if cl.IntersectsWith(cx, cy, scale) {
		if cl.Scale < scale-span {
			for _, d := range cl.Data {
				d.Export(dm)
			}
		} else {
			cl.EachChildren(func(dx int, dy int, c *Cluster) {
				if c != nil {
					c.ViewMap(dm, cx, cy, scale, span)
				}
			})
		}
	}
}

func (cl *Cluster) EachChildren(callback func(int, int, *Cluster)) {
	if cl.Scale > Const.MinScale {
		for _, dy := range []int{-1, +1} {
			for _, dx := range []int{-1, +1} {
				callback(dx, dy, cl.FindChild(dx, dy))
			}
		}
	}
}

func (cl *Cluster) BeforeDelete() {
	if cl.Parent != nil {
		cl.Parent.UnResolve(cl)
	}
}

func (cl *Cluster) UnResolve(args ...interface{}) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Cluster:
			for _, list := range cl.Children {
				for _, child := range list {
					if child == cl {
						child = nil
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

func (cl *Cluster) IntersectsWith(cx float64, cy float64, scale float64) bool {
	myL := math.Pow(2, cl.Scale) / 2
	othL := math.Pow(2, scale) / 2

	return math.Max(cl.X-myL, cx-othL) <= math.Min(cl.X+myL, cx+othL) &&
		math.Max(cl.Y-myL, cy-othL) <= math.Min(cl.Y+myL, cy+othL)
}
