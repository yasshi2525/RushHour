package entities

import "math"

// Point represents geographical location on game map
type Point struct {
	X float64 `gorm:"index;not null"`
	Y float64 `gorm:"index;not null"`
}

// Dist calculate a distance between two Point
func (p *Point) Dist(oth *Point) float64 {
	return math.Sqrt((oth.X-p.X)*(oth.X-p.X) + (oth.Y-p.Y)*(oth.Y-p.Y))
}

// IsIn returns true when Point is in specified area
func (p *Point) IsIn(center *Point, scale float64) bool {
	len := math.Pow(2, scale)

	return p.X > center.X-len/2 &&
		p.X < center.X+len/2 &&
		p.Y > center.Y-len/2 &&
		p.Y < center.Y+len/2
}

// Junction is a logical Point that connects Edges.
// There is more than two Junction on same geographically xy if Human cannot move.
type Junction struct {
	Point
	Out []*Step `gorm:"-"`
	In  []*Step `gorm:"-"`
}

// NewJunction create Juntion
func NewJunction(x float64, y float64) Junction {
	return Junction{
		Point: Point{X: x, Y: y},
		Out:   []*Step{},
		In:    []*Step{},
	}
}

// Step represents two Junction is logically connected.
// Step is out of target for persistence because it can intruduce by other resources.
type Step struct {
	ID     uint
	From   *Junction
	To     *Junction
	Weight float64
}

// Cost is calculated by distance * weight of Step
func (s *Step) Cost() float64 {
	return s.From.Dist(&s.To.Point) * s.Weight
}
