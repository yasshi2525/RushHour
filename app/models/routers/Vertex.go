package routers

type Vertex struct {
	x   float64
	y   float64
	out []Edge
	in  []Edge
}

func NewVertex(x float64, y float64) *Vertex {
	inst := new(Vertex)
	inst.x = x
	inst.y = y
	return inst
}

func (v *Vertex) Dist(oth *Vertex) float64 {
	return 0
}

func (v *Vertex) isIn(center *Vertex, scale float64) bool {
	return false
}
