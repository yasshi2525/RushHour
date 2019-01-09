package entities

import (
	"fmt"
)

// Agent is the wrapper of Human.
// Agent can move concerting minimum cost route.
type Agent struct {
	Human   *Human
	Current *Step
}

// NewAgent creates instance
func (m *Model) NewAgent(h *Human) *Agent {
	a := &Agent{
		Human: h,
	}
	a.Init()
	m.Add(a)
	return a
}

// Init do nothing
func (a *Agent) Init() {

}

// Idx returns unique id field.
func (a *Agent) Idx() uint {
	return a.Human.ID
}

// Type returns type of entitiy
func (a *Agent) Type() ModelType {
	return AGENT
}

// Consume makes Human moves
func (a *Agent) Consume(interval float64) {
	//TODO
}

// String represents status
func (a *Agent) String() string {
	idstr := "?"
	posstr := ""
	if a.Human != nil {
		idstr = fmt.Sprintf("%d", a.Human.ID)
		posstr = fmt.Sprintf(":%s", a.Human.Pos())
	}
	return fmt.Sprintf("%s(%s):%v%s", a.Type().Short(),
		idstr, a.Current, posstr)
}
