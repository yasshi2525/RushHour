package v1

import (
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
	return c.RenderJSON(genResponse(true, entities.JsonPlayer(services.Model.Players)))
}

// Diff returns only diff
func (c APIv1Game) Diff() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	return c.RenderJSON(
		genResponse(
			true,
			services.ViewMap(500, 500, 10, time.Now().Add(time.Duration(-1)*time.Minute))))
}

// Departure returns result of rail node creation
func (c APIv1Game) Departure() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	p := &PointRequest{}
	if errs := p.Parse(c.Params.Form); len(errs) > 0 {
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
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	p := &PointRequest{}
	errs := p.Parse(c.Params.Form)
	e, err := validateEntity(entities.RAILNODE, c.Params.Form.Get("rnid"))
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
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()
	o := &OwnerRequest{}
	errs := o.Parse(c.Params.Form)
	e1, err1 := validateEntity(entities.RAILNODE, c.Params.Form.Get("from"))
	e2, err2 := validateEntity(entities.RAILNODE, c.Params.Form.Get("to"))
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
	re, err := services.ConnectRailNode(o.O, from, to, o.Scale)
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}
	return c.RenderJSON(genResponse(true, &struct {
		In  *entities.DelegateRailEdge `json:"e1"`
		Out *entities.DelegateRailEdge `json:"e2"`
	}{re, re.Reverse}))
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
