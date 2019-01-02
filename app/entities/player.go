package entities

// PlayerType represents authenticate level
type PlayerType uint

// PlayerType represents authenticate level
const (
	Admin PlayerType = iota + 1
	Normal
	Guest
)

// Player represents user information
type Player struct {
	Model

	Level       PlayerType `gorm:"not null"`
	DisplayName string     `gorm:"not null"`
	LoginID     string     `gorm:"not null;index" json:"-"`
	Password    string     `gorm:"not null" json:"-"`
}

// NewPlayer create instance
func NewPlayer(id uint) *Player {
	return &Player{
		Model: NewModel(id),
	}
}

// Idx returns unique id field.
func (o *Player) Idx() uint {
	return o.ID
}

// Init do nothing
func (o *Player) Init() {
	o.Model.Init()
}

// ResolveRef do nothing for implementing Resolvable
func (o *Player) ResolveRef() {
	// do-nothing
}
