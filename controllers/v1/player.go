package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/entities"
	"github.com/yasshi2525/RushHour/services"
)

// Players returns list of player
// @Description list of player
// @Tags []entities.Player
// @Summary list of player
// @Accept json
// @Produce json
// @Success 200 {array} entities.Player "list of player"
// @Failure 503 {object} errInfo "under maintenance"
// @Router /players [get]
func Players(c *gin.Context) {
	c.Set(keyOk, entities.JSONPlayer(services.Model.Players))
}
