package entities

type Transport struct {
	FromPlatform *Platform
	ToPlatform   *Platform
	Via          *LineTask
	Value        float64
}

func (ts *Transport) ExportStep(id uint) *Step {
	s := NewStep(id, ts.FromPlatform, ts.ToPlatform)
	s.By = ts.Via
	s.Transport = ts.Value
	return s
}

type Track struct {
	FromNode *RailNode
	ToNode   *RailNode
	Via      *RailEdge
	Value    float64
}
