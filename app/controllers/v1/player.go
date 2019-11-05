package v1

import (
	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

// Players returns list of player
// @Description list of player
// @Tags []entities.Player
// @Summary list of player
// @Accept  query
// @Produce  json
// @Success 200 {object} []entities.Player "list of player"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /players [get]
func Players(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	c.Set(keyOk, entities.JSONPlayer(services.Model.Players))
}
