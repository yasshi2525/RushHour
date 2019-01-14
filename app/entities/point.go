package entities

import (
	"fmt"
	"math"
	"math/rand"
)

// Point represents geographical location on game map
type Point struct {
	X float64 `gorm:"index;not null" json:"x"`
	Y float64 `gorm:"index;not null" json:"y"`
}

// NewPoint create Point
func NewPoint(x float64, y float64) Point {
	return Point{X: x, Y: y}
}

// Pos returns self
func (p *Point) Pos() *Point {
	return p
}

// IsIn returns true when Point is in specified area
func (p *Point) IsIn(x float64, y float64, scale float64) bool {
	len := math.Pow(2, scale)

	return p.X > x-len/2 && p.X < x+len/2 && p.Y > y-len/2 && p.Y < y+len/2
}

// IsInLine returns true when this or to or center is in.
func (p *Point) IsInLine(to Locationable, x float64, y float64, scale float64) bool {
	return p.IsIn(x, y, scale) ||
		p.Center(to).IsIn(x, y, scale) ||
		to.Pos().IsIn(x, y, scale)
}

// Dist calculate a distance between two Point
func (p *Point) Dist(oth Locationable) float64 {
	return math.Sqrt((oth.Pos().X-p.X)*(oth.Pos().X-p.X) + (oth.Pos().Y-p.Y)*(oth.Pos().Y-p.Y))
}

// Center returns devided point.
func (p *Point) Center(to Locationable) *Point {
	return p.Div(to, 0.5)
}

// Div returns dividing point to certain ratio.
func (p *Point) Div(to Locationable, progress float64) *Point {
	return &Point{
		X: p.X*progress + to.Pos().X*(1-progress),
		Y: p.Y*progress + to.Pos().Y*(1-progress),
	}
}

func (p *Point) Rand(max float64) *Point {
	dist := rand.Float64() * max
	rad := rand.Float64() * math.Pi * 2
	return &Point{
		X: p.X + dist*math.Cos(rad),
		Y: p.Y + dist*math.Sin(rad),
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", p.X, p.Y)
}
