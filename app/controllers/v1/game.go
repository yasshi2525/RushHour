package v1

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

type gameStatus struct {
	Status bool `json:"status"`
}

// GameStatus returns game status
// @Description game status
// @Tags gameStatus
// @Summary game status
// @Produce json
// @Success 200 {object} gameStatus "game status"
// @Router /game [get]
func GameStatus(c *gin.Context) {
	c.Set(keyOk, &gameStatus{services.IsInOperation()})
}

// StartGame returns result of game starting
// @Description result of game starting
// @Tags gameStatus
// @Summary start game
// @Produce json
// @Success 200 {object} gameStatus "game status"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /game/start [post]
func StartGame(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	if o.Level != entities.Admin {
		c.Set(keyErr, fmt.Errorf("permission denied"))
	} else {
		if !services.IsInOperation() {
			services.Start()
		}
		c.Set(keyOk, &gameStatus{services.IsInOperation()})
	}
}

// StopGame returns result of game stopping
// @Description result of game stopping
// @Tags gameStatus
// @Summary stop game
// @Produce json
// @Success 200 {object} gameStatus "game status"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /game/stop [post]
func StopGame(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	if o.Level != entities.Admin {
		c.Set(keyErr, fmt.Errorf("permission denied"))
	} else {
		if services.IsInOperation() {
			services.Stop()
		}
		c.Set(keyOk, &gameStatus{services.IsInOperation()})
	}
}

type purgeStatus struct {
	Purge bool `json:"purge"`
}

// PurgeUserData deletes all user data
// @Description result of purging
// @Tags gameStatus
// @Summary start game
// @Produce json
// @Success 200 {object} gameStatus "game status"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /game/start [post]
func PurgeUserData(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	if o.Level != entities.Admin {
		c.Set(keyErr, fmt.Errorf("permission denied"))
	} else if err := services.Purge(o.O); err != nil {
		c.Set(keyErr, err)
	} else {
		c.Set(keyOk, &purgeStatus{true})
	}
}
