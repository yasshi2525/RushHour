package entities

import (
	"math"
)

// Shape containts geometorical information.
type Shape struct {
	P1       *Point   `gorm:"-" json:"-"`
	P2       *Point   `gorm:"-" json:"-"`
	Children []*Shape `gorm:"-" json:"-"`
	Referred []*Shape `gorm:"-" json:"-"`
}

// NewShapeGroup create instance that groups other Shapes.
func NewShapeGroup() Shape {
	return Shape{nil, nil, []*Shape{}, []*Shape{}}
}

// NewShapeNode create instance that represents Point.
func NewShapeNode(pos *Point) Shape {
	return Shape{pos, nil, nil, []*Shape{}}
}

// NewShapeEdge create instance that represents Edge, which consists in two Point.
func NewShapeEdge(from *Point, to *Point) Shape {
	return Shape{from, to, nil, []*Shape{}}
}

// Append add specified child to spape's list.
func (sh *Shape) Append(child *Shape) {
	sh.Children = append(sh.Children, child)
}

// Refer add specified child to referred list.
func (sh *Shape) Refer(ref *Shape) {
	sh.Referred = append(sh.Referred, ref)
}

// Delete removes specified child from shape's list.
func (sh *Shape) Delete(child *Shape) {
	for i, c := range sh.Children {
		if c == child {
			buf := append(sh.Children[:i], sh.Children[i+1:]...)
			sh.Children = make([]*Shape, len(buf))
			copy(sh.Children, buf)
			return
		}
	}
}

// UnRefer removes specified child from referred's list.
func (sh *Shape) UnRefer(ref *Shape) {
	for i, r := range sh.Referred {
		if r == ref {
			buf := append(sh.Referred[:i], sh.Referred[i+1:]...)
			sh.Referred = make([]*Shape, len(buf))
			copy(sh.Referred, buf)
			return
		}
	}
}

// Pos returns location.
func (sh *Shape) Pos() *Point {
	if sh.Children != nil {
		length := len(sh.Children)
		if length == 0 {
			return nil
		}
		var sumX, sumY float64
		for _, c := range sh.Children {
			x, y := c.Pos().Flat()
			sumX += x
			sumY += y
		}
		return &Point{sumX / float64(length), sumY / float64(length)}
	}
	if sh.P1 == nil {
		return &Point{}
	}
	if sh.P2 == nil {
		return sh.P1
	}
	return sh.P1.Center(sh.P2)
}

// IsIn returns it should be view or not.
func (sh *Shape) IsIn(x float64, y float64, scale float64) bool {
	for _, r := range sh.Referred {
		if r.IsIn(x, y, scale) {
			return true
		}
	}
	if sh.Children != nil {
		for _, c := range sh.Children {
			if c.IsIn(x, y, scale) {
				return true
			}
		}
		return false
	}
	if sh.P1 == nil {
		return false
	}
	if sh.P2 == nil {
		return sh.P1.IsIn(x, y, scale)
	}
	return sh.P1.IsInLine(sh.P2, x, y, scale)
}

// Dist returns length of this Edge.
func (sh *Shape) Dist() float64 {
	if sh.P1 == nil || sh.P2 == nil {
		return 0
	}
	return sh.P1.Dist(sh.P2)
}

// Angle returns angle with 'to' object.
func (sh *Shape) Angle(to *Shape) float64 {
	if sh.P1 == nil || sh.P2 == nil {
		return 0
	}
	v := sh.P2.Sub(sh.P1).Unit()
	u := to.P2.Sub(to.P1).Unit()
	theta := math.Acos(v.InnerProduct(u))
	if math.IsNaN(theta) {
		return math.Pi
	}
	return theta
}

// Div returns Point which divides 'prog' ratio to Edge.
func (sh *Shape) Div(prog float64) *Point {
	if sh.P1 == nil {
		return &Point{}
	}
	if sh.P2 == nil {
		return sh.P1
	}
	return sh.P1.Div(sh.P2, prog)
}
