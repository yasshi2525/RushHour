package services

import (
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
)

type gameMap struct {
	Players    []*entities.Player    `json:"players"`
	Residences []*entities.Residence `json:"residences"`
	Companies  []*entities.Company   `json:"companies"`
	RailNodes  []*entities.RailNode  `json:"rail_nodes"`
}

// ViewMap immitates user requests view
// TODO remove
func ViewMap() interface{} {
	start := time.Now()
	defer WarnLongExec(start, 2, "ユーザ表示要求", false)

	game := &gameMap{
		Players:    []*entities.Player{},
		Residences: []*entities.Residence{},
		Companies:  []*entities.Company{},
		RailNodes:  []*entities.RailNode{},
	}

	MuStatic.RLock()
	defer MuStatic.RUnlock()

	for _, val := range Static.Players {
		game.Players = append(game.Players, val)
	}
	for _, val := range Static.Residences {
		game.Residences = append(game.Residences, val)
	}
	for _, val := range Static.Companies {
		game.Companies = append(game.Companies, val)
	}
	for _, val := range Static.RailNodes {
		game.RailNodes = append(game.RailNodes, val)
	}

	return game
}
