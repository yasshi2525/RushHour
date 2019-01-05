package entities

type Track struct {
	FromPlatform *Platform
	ToPlatform   *Platform
	Via          *LineTask
	Value        float64
}

func NewTrack(from *Platform, to *Platform, via *LineTask, v float64) *Track {
	return &Track{from, to, via, v}
}

func (tr *Track) ExportStep(id uint) *Step {
	s := NewStep(id, tr.FromPlatform, tr.ToPlatform)
	s.By = tr.Via
	s.Transport = tr.Value
	return s
}
