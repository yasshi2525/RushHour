package entities

import "github.com/jinzhu/gorm"

// Player represents user information
type Player struct {
	gorm.Model

	DisplayName string `gorm:"not null"`
	LoginID     string `gorm:"not null,index"`
	Password    string `gorm:"not null"`
}

// ResolveRef do nothing for implementing Resolvable
func (p *Player) ResolveRef() {
	// do-nothing
}
