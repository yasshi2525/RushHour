package services

import (
	"fmt"
	"strings"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string) (*entities.Player, error) {
	// duplication check
	for _, oth := range Repo.Static.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			revel.AppLog.Warnf("Login ID \"%s\" is already used.", loginid)
			return nil, fmt.Errorf("login ID \"%s\" is already used", loginid)
		}
	}

	player := entities.NewPlayer(GenID(entities.PLAYER))
	player.LoginID = loginid
	player.DisplayName = displayname
	player.Password = password

	Repo.Static.Players[player.ID] = player
	revel.AppLog.Infof("Player(%d) %s was created", player.ID, loginid)

	return player, nil
}

// FetchOwner fetch loginid Player.
func FetchOwner(loginid string) (*entities.Player, error) {
	for _, oth := range Repo.Static.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			return oth, nil
		}
	}
	return nil, fmt.Errorf("login ID \"%s\" was not found", loginid)
}
