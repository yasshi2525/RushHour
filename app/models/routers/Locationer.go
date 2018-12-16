package routers

type Locationer interface {
	Get() (float64, float64)
	Dist(oth Locationer) float64
	IsIn(x float64, y float64) bool
}
