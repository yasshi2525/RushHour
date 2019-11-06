package v1

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
)

type deptResponse struct {
	RailNode *entities.DelegateRailNode `json:"rn"`
}

// Depart returns result of rail node creation
// @Description result of rail node creation
// @Tags deptResponse
// @Summary depart
// @Accept json
// @Produce json
// @Param x body number true "x coordinate"
// @Param y body number true "y coordinate"
// @Param scale body number true "width,height(100%)=2^scale"
// @Success 200 {object} deptResponse "created rail node"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /rail_nodes [post]
func Depart(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	params := pointRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else {
		x, y, sc := params.export()
		if rn, err := services.CreateRailNode(o, x, y, sc); err != nil {
			c.Set(keyErr, err)
		} else {
			c.Set(keyOk, &deptResponse{RailNode: rn})
		}
	}
}

type extendRequest struct {
	pointRequest
	RailNode string `form:"rnid" json:"rnid" validate:"required,numeric"`
}

type extendResponse struct {
	RailNode *entities.DelegateRailNode `json:"rn"`
	In       *entities.DelegateRailEdge `json:"e1"`
	Out      *entities.DelegateRailEdge `json:"e2"`
}

// Extend returns result of rail node extension
// @Description result of rail node extension
// @Tags extendResponse
// @Summary extend
// @Accept json
// @Produce json
// @Param x body number true "x coordinate"
// @Param y body number true "y coordinate"
// @Param scale body number true "width,height(100%)=2^scale"
// @Param rnid body integer true "tail rail node id"
// @Success 200 {object} extendResponse "extend rail node"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /rail_nodes/extend [post]
func Extend(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	params := extendRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else if rn, err := validateEntity(entities.RAILNODE, params.RailNode); err != nil {
		c.Set(keyErr, err)
	} else {
		x, y, sc := params.export()
		if to, re, err := services.ExtendRailNode(o, rn.(*entities.RailNode), x, y, sc); err != nil {
			c.Set(keyErr, err)
		} else {
			c.Set(keyOk, &extendResponse{to, re, re.Reverse})
		}
	}
}

type connectRequest struct {
	scaleRequest
	From string `form:"from" json:"from" validate:"required,numeric"`
	To   string `form:"to" json:"to" validate:"required,numeric"`
}

type connectResponse struct {
	In  *entities.DelegateRailEdge `json:"e1"`
	Out *entities.DelegateRailEdge `json:"e2"`
}

// Connect returns result of rail connection
// @Description result of rail node connection
// @Tags connectResponse
// @Summary connect
// @Accept json
// @Produce json
// @Param x body number true "x coordinate"
// @Param y body number true "y coordinate"
// @Param scale body number true "width,height(100%)=2^scale"
// @Param from body integer true "from rail node id"
// @Param to body integer true "to rail node id"
// @Success 200 {object} connectResponse "connect rail node"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /rail_nodes/connect [post]
func Connect(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	params := connectRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else if from, err := validateEntity(entities.RAILNODE, params.From); err != nil {
		c.Set(keyErr, err)
	} else if to, err := validateEntity(entities.RAILNODE, params.To); err != nil {
		c.Set(keyErr, err)
	} else {
		sc := params.export()
		if re, err := services.ConnectRailNode(o, from.(*entities.RailNode), to.(*entities.RailNode), sc); err != nil {
			c.Set(keyErr, err)
		} else {
			c.Set(keyOk, &connectResponse{re, re.Reverse})
		}
	}
}

type removeRailNodeRequest struct {
	RailNode string `form:"rnid" json:"rnid" validate:"required,numeric"`
}

func (v *removeRailNodeRequest) export() uint {
	rnid, _ := strconv.ParseFloat(v.RailNode, 64)
	return uint(rnid)
}

type removeRailNodeResponse struct {
	RailNode string `json:"rnid"`
}

// RemoveRailNode returns result of rail deletion
// @Description result of rail node deletion
// @Tags removeRailNodeResponse
// @Summary remove rail node
// @Accept json
// @Produce json
// @Param rnid body integer true "rail node id"
// @Success 200 {object} removeRailNodeResponse "connect rail node"
// @Success 400 {object} errInfo "reason of fail"
// @Failure 401 {object} errInfo "invalid jwt"
// @Router /rail_nodes [delete]
func RemoveRailNode(c *gin.Context) {
	o := authorize(c)
	if o == nil {
		return
	}
	params := removeRailNodeRequest{}
	if err := c.Bind(&params); err != nil {
		c.Set(keyErr, err)
	} else if err := services.RemoveRailNode(o, params.export()); err != nil {
		c.Set(keyErr, err)
	} else {
		c.Set(keyOk, &removeRailNodeResponse{params.RailNode})
	}
}
