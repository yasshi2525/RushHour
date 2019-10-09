package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"

	"github.com/revel/revel"
	"github.com/yasshi2525/RushHour/app/services"
)

// APIv1Game is controller for REST API
type APIv1Game struct {
	*revel.Controller
}

// Index returns gamemap
func (c APIv1Game) Index() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	var params = make(map[string]float64)

	for _, p := range []string{"cx", "cy", "scale", "delegate"} {
		if v, err := strconv.ParseFloat(c.Params.Get(p), 64); err == nil {
			params[p] = v
		}
	}

	return c.RenderJSON(
		genResponse(true, services.ViewDelegateMap(params["cx"], params["cy"], params["scale"], params["delegate"])))
}

// Players returns list of player
func (c APIv1Game) Players() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	return c.RenderJSON(genResponse(true, entities.JSONPlayer(services.Model.Players)))
}

// Departure returns result of rail node creation
func (c APIv1Game) Departure() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	p := &PointRequest{}
	if errs := p.Parse(token, json); len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	rn, err := services.CreateRailNode(p.O, p.X, p.Y, p.Scale)
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	return c.RenderJSON(genResponse(true, &struct {
		RailNode *entities.DelegateRailNode `json:"rn"`
	}{rn}))
}

// Extend returns result of rail node extension
func (c APIv1Game) Extend() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	p := &PointRequest{}
	errs := p.Parse(token, json)
	e, err := validateEntity(entities.RAILNODE, json["rnid"])
	if err != nil {
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	from := e.(*entities.RailNode)
	to, re, err := services.ExtendRailNode(p.O, from, p.X, p.Y, p.Scale)
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	return c.RenderJSON(genResponse(true, &struct {
		RailNode *entities.DelegateRailNode `json:"rn"`
		In       *entities.DelegateRailEdge `json:"e1"`
		Out      *entities.DelegateRailEdge `json:"e2"`
	}{to, re, re.Reverse}))
}

// Connect returns result of rail connection
func (c APIv1Game) Connect() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	s := &ScaleRequest{}
	errs := s.Parse(token, json)
	e1, err1 := validateEntity(entities.RAILNODE, json["from"])
	e2, err2 := validateEntity(entities.RAILNODE, json["to"])
	if err1 != nil {
		errs = append(errs, err1.Error())
	}
	if err2 != nil {
		errs = append(errs, err2.Error())
	}
	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	from := e1.(*entities.RailNode)
	to := e2.(*entities.RailNode)
	re, err := services.ConnectRailNode(s.O, from, to, s.Scale)
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	return c.RenderJSON(genResponse(true, &struct {
		In  *entities.DelegateRailEdge `json:"e1"`
		Out *entities.DelegateRailEdge `json:"e2"`
	}{re, re.Reverse}))
}

// RemoveRailNode returns result of rail deletion.
func (c APIv1Game) RemoveRailNode() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	o := &OwnerRequest{}
	errs := o.Parse(token, json)

	id, ok := json["id"].(float64)
	if !ok {
		errs = append(errs, fmt.Sprintf("parse id failed: %v", json["id"]))
	}
	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}
	if err := services.RemoveRailNode(o.O, uint(id)); err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}
	return c.RenderJSON(genResponse(true, &struct {
		DeleteID uint `json:"id"`
	}{uint(id)}))
}

func genResponse(status bool, results interface{}) interface{} {
	return &struct {
		Status    bool        `json:"status"`
		Timestamp int64       `json:"timestamp"`
		Results   interface{} `json:"results"`
	}{
		Status:    true,
		Timestamp: time.Now().Unix(),
		Results:   results,
	}
}

// GetToken returns auth token from session.
func (c APIv1Game) getToken() (string, error) {
	if token, err := c.Session.Get("token"); err != nil {
		return "", err
	} else {
		return token.(string), nil
	}
}
