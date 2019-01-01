package entities

import (
	"fmt"
	"math"
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
	return p.Div(to, 0.5)
}

// Div returns dividing point to certain ratio.
func (p *Point) Div(to *Point, progress float64) *Point {
	return &Point{
		X: p.X*progress + to.X*(1-progress),
		Y: p.Y*progress + to.Y*(1-progress),
	}
}

func (p *Point) String() string {
	return fmt.Sprintf("(%f, %f)", p.X, p.Y)
}
