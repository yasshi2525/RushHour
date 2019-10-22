package v1

import (
	"fmt"
	"strconv"
	"time"

	"github.com/yasshi2525/RushHour/app/entities"
	"github.com/yasshi2525/RushHour/app/services"
	"github.com/yasshi2525/RushHour/app/services/auth"

	"github.com/revel/revel"
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

// Login returns result of password login
func (c APIv1Game) Login() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	id, okid := json["id"].(string)
	password, okpw := json["password"].(string)

	if !okid || !okpw {
		return c.RenderJSON(genResponse(false, "id or password is invalid"))
	}
	if o, err := services.PasswordSignIn(id, password); err != nil {
		return c.RenderJSON(genResponse(false, err))
	} else {
		c.Session.Set("token", o.Token)
		return c.RenderJSON(genResponse(true, o))
	}
}

// Register returns result of password sign up
func (c APIv1Game) Register() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	id, okid := json["id"].(string)
	password, okpw := json["password"].(string)
	name, oknm := json["name"].(string)
	hue, okh := json["hue"].(float64)

	errs := []error{}

	if !okid {
		errs = append(errs, fmt.Errorf("id is invalid"))
	}
	if !okpw {
		errs = append(errs, fmt.Errorf("password is invalid"))
	}
	if !oknm {
		errs = append(errs, fmt.Errorf("name is invalid"))
	}
	if !okh {
		errs = append(errs, fmt.Errorf("hue is invalid"))
	}
	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}
	if o, err := services.PasswordSignUp(id, name, password, int(hue), entities.Normal); err != nil {
		return c.RenderJSON(genResponse(false, err))
	} else {
		c.Session.Set("token", o.Token)
		return c.RenderJSON(genResponse(true, o))
	}
}

// GetSettings returns the list of customizable attributes
func (c APIv1Game) GetSettings() revel.Result {
	services.MuModel.RLock()
	defer services.MuModel.RUnlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	o := &OwnerRequest{}
	errs := o.Parse(token)

	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	return c.RenderJSON(genResponse(true, services.GetAccountSettings(o.O)))
}

// ChangeSettings returns the result of change profile
func (c APIv1Game) ChangeSettings() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	o := &OwnerRequest{}
	errs := o.Parse(token)

	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	json := make(map[string]interface{})
	c.Params.BindJSON(&json)

	var value interface{}
	var ok bool

	if value, ok = json["value"]; !ok {
		return c.RenderJSON(genResponse(false, []error{fmt.Errorf("value is not set")}))
	}
	res := c.Params.Get("resname")
	switch res {
	case "custom_name":
		if v, ok := value.(string); !ok {
			return c.RenderJSON(genResponse(false, []error{fmt.Errorf("%s is not string: %s", res, value)}))
		} else {
			o.O.CustomDisplayName = auth.Encrypt(v)
		}
	case "use_cname":
		if v, ok := value.(bool); !ok {
			return c.RenderJSON(genResponse(false, []error{fmt.Errorf("%s is not bool: %s", res, value)}))
		} else {
			o.O.UseCustomDisplayName = v
		}
	default:
		return c.RenderJSON(genResponse(false, []error{fmt.Errorf("invalid attribute %s", res)}))
	}
	return c.RenderJSON(genResponse(true, &struct {
		Player *entities.Player `json:"my"`
		Key    string           `json:"key"`
		Value  interface{}      `json:"value"`
	}{
		Player: o.O,
		Key:    res,
		Value:  value,
	}))
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
	errs := o.Parse(token)

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

func (c APIv1Game) GameStatus() revel.Result {
	return c.RenderJSON(genResponse(true, services.IsInOperation()))
}

func (c APIv1Game) StartGame() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	o := &OwnerRequest{}
	errs := o.Parse(token)

	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	if o.O.Level != entities.Admin {
		return c.RenderJSON(genResponse(false, []error{fmt.Errorf("permission denied")}))
	}

	if !services.IsInOperation() {
		services.Start()
	}
	return c.RenderJSON(genResponse(true, true))
}

func (c APIv1Game) StopGame() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	o := &OwnerRequest{}
	errs := o.Parse(token)

	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	if o.O.Level != entities.Admin {
		return c.RenderJSON(genResponse(false, []error{fmt.Errorf("permission denied")}))
	}

	if services.IsInOperation() {
		services.Stop()
	}
	return c.RenderJSON(genResponse(true, false))
}

// PurgeUserData deletes all user data.
func (c APIv1Game) PurgeUserData() revel.Result {
	services.MuModel.Lock()
	defer services.MuModel.Unlock()

	token, err := c.getToken()
	if err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}

	o := &OwnerRequest{}
	errs := o.Parse(token)

	if len(errs) > 0 {
		return c.RenderJSON(genResponse(false, errs))
	}

	if o.O.Level != entities.Admin {
		return c.RenderJSON(genResponse(false, []error{fmt.Errorf("permission denied")}))
	}

	if err := services.Purge(o.O); err != nil {
		return c.RenderJSON(genResponse(false, []error{err}))
	}
	return c.RenderJSON(genResponse(true, true))
}

func genResponse(status bool, results interface{}) interface{} {
	var details interface{}

	switch obj := results.(type) {
	case []error:
		outputs := []string{}
		for _, err := range obj {
			outputs = append(outputs, err.Error())
		}
		details = outputs
	case error:
		details = []string{obj.Error()}
	default:
		details = results
	}

	return &struct {
		Status    bool        `json:"status"`
		Timestamp int64       `json:"timestamp"`
		Results   interface{} `json:"results"`
	}{
		Status:    status,
		Timestamp: time.Now().Unix(),
		Results:   details,
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
