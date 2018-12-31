package services

import (
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/revel/revel"

	"github.com/yasshi2525/RushHour/app/entities"
)

// CreatePlayer creates player.
func CreatePlayer(loginid string, displayname string, password string) (*entities.Player, error) {
	id := uint(atomic.AddUint64(NextID.Static[PLAYER], 1))

	// duplication check
	for _, oth := range Static.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			revel.AppLog.Warnf("Login ID \"%s\" is already used.", loginid)
			return nil, fmt.Errorf("login ID \"%s\" is already used", loginid)
		}
	}

	player := &entities.Player{
		LoginID:     loginid,
		DisplayName: displayname,
		Password:    password,
	}

	Static.Players[id] = player
	revel.AppLog.Infof("Player(%d) %s was created", id, loginid)

	return player, nil
}

// FetchOwner fetch loginid Player.
func FetchOwner(loginid string) (*entities.Player, error) {
	for _, oth := range Static.Players {
		if strings.Compare(loginid, oth.LoginID) == 0 {
			return oth, nil
		}
	}
	return nil, fmt.Errorf("login ID \"%s\" was not found", loginid)
}
