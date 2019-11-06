package entities

import (
	"fmt"
	"math"
	"math/rand"
)

// Point represents geographical location on game map
type Point struct {
	X float64 `gorm:"index" json:"x"`
	Y float64 `gorm:"index" json:"y"`
}

// NewPoint create Point
func NewPoint(x float64, y float64) Point {
	return Point{X: x, Y: y}
}

// IsIn returns true when Point is in specified area
func (p *Point) IsIn(x float64, y float64, scale float64) bool {
	len := math.Pow(2, scale)

	return p.X >= x-len/2 && p.X < x+len/2 && p.Y >= y-len/2 && p.Y < y+len/2
}

// IsInLine returns true when this or to or center is in.
func (p *Point) IsInLine(to *Point, x float64, y float64, scale float64) bool {
	return p.IsIn(x, y, scale) ||
		p.Center(to).IsIn(x, y, scale) ||
		to.IsIn(x, y, scale)
}

// Dist calculate a distance between two Point
func (p *Point) Dist(oth *Point) float64 {
	return math.Sqrt((oth.X-p.X)*(oth.X-p.X) + (oth.Y-p.Y)*(oth.Y-p.Y))
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

// Rand generates other Point randaomly within 'max' distance.
func (p *Point) Rand(max float64) *Point {
	dist := rand.Float64() * max
	rad := rand.Float64() * math.Pi * 2
	return &Point{
		X: p.X + dist*math.Cos(rad),
		Y: p.Y + dist*math.Sin(rad),
	}
}

// Flat returns position as two value
func (p *Point) Flat() (float64, float64) {
	return p.X, p.Y
}

// Sub returns new Point which is substracted by 'to' object
func (p *Point) Sub(to *Point) *Point {
	return &Point{p.X - to.X, p.Y - to.Y}
}

// Unit returns unit vector of this Point origined by (0, 0)
func (p *Point) Unit() *Point {
	length := p.Dist(&Point{})
	return &Point{p.X / length, p.Y / length}
}

// InnerProduct returns inner product with 'to' object.
func (p *Point) InnerProduct(to *Point) float64 {
	return p.X*to.X + p.Y*to.Y
}

// Clone returns same value but referrence is different value object.
func (p *Point) Clone() *Point {
	return &Point{p.X, p.Y}
}

func (p Point) String() string {
	return fmt.Sprintf("(%.2f,%.2f)", p.X, p.Y)
}
