package services

import (
	"fmt"
	"strings"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string, level entities.PlayerType) (*entities.Player, error) {
	// duplication check
	for _, oth := range Model.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			revel.AppLog.Warnf("Login ID \"%s\" is already used.", loginid)
			return nil, fmt.Errorf("login ID \"%s\" is already used", loginid)
		}
	}

	player := Model.NewPlayer()
	player.LoginID = loginid
	player.DisplayName = displayname
	player.Password = password
	player.Level = level

	return player, nil
}

// FetchOwner fetch loginid Player.
func FetchOwner(loginid string) (*entities.Player, error) {
	for _, oth := range Model.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			return oth, nil
		}
	}
	return nil, fmt.Errorf("login ID \"%s\" was not found", loginid)
}
