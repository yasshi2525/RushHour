package entities

// Player represents user information
type Player struct {
	Model

	DisplayName string `gorm:"not null"`
	LoginID     string `gorm:"not null,index json:"-""`
	Password    string `gorm:"not null" json:"-"`
}

// ResolveRef do nothing for implementing Resolvable
func (p *Player) ResolveRef() {
	// do-nothing
}
