package entities

import "math"

// Point represents geographical location on game map
type Point struct {
	X float64 `gorm:"index;not null" json:"x"`
	Y float64 `gorm:"index;not null" json:"y"`
}

// NewPoint create Point
func NewPoint(x float64, y float64) *Point {
	return &Point{X: x, Y: y}
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

// Center returns devided point.
func (p *Point) Center(to *Point) *Point {
	return &Point{
		X: (p.X + to.X) / 2,
		Y: (p.Y + to.Y) / 2,
	}
}

// Junction is a logical Point that connects Edges.
// There is more than two Junction on same geographically xy if Human cannot move.
type Junction interface {
	Pos() *Point
	OutStep() map[uint]*Step
	InStep() map[uint]*Step
}

// Step represents two Junction is logically connected.
// Step is out of target for persistence because it can intruduce by other resources.
type Step struct {
	ID     uint
	From   Junction
	To     Junction
	Weight float64
}

// NewStep create new instance and relation to Junction
func NewStep(id uint, from Junction, to Junction, weight float64) *Step {
	step := &Step{
		ID:     id,
		From:   from,
		To:     to,
		Weight: weight,
	}
	from.OutStep()[step.ID] = step
	to.InStep()[step.ID] = step

	return step
}

// Cost is calculated by distance * weight of Step.
func (s *Step) Cost() float64 {
	return s.From.Pos().Dist(s.To.Pos()) * s.Weight
}

// Unrelate delete selt from related Junction.
func (s *Step) Unrelate() {
	delete(s.From.OutStep(), s.ID)
	delete(s.To.InStep(), s.ID)
}
