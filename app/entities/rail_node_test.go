package entities

import "testing"

func TestNewRailNode(t *testing.T) {
	m := NewModel()
	o := m.NewPlayer()

	var x, y float64 = 10.0, 20.0

	rn := m.NewRailNode(o, x, y)

	TestCases{
		{"O", rn.O, o},
		{"X", rn.X, x},
		{"Y", rn.Y, y},
		{"model", m.RailNodes[rn.Idx()], rn},
	}.Assert(t)
}

func TestExtend(t *testing.T) {
	t.Run("without line", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		from := m.NewRailNode(o, 0, 0)
		var x, y float64 = 10.0, 20.0
		to, e1 := from.Extend(x, y)

		TestCases{
			{"to.X", to.X, x},
			{"to.Y", to.Y, y},
			{"e1.from", e1.FromNode, from},
			{"e1.to", e1.ToNode, to},
			{"e1.reverse", e1.Reverse.Reverse, e1},
			{"e2.from", e1.Reverse.FromNode, to},
			{"e2.to", e1.Reverse.ToNode, from},
		}.Assert(t)
	})
	t.Run("with autoExt", func(t *testing.T) {
		m := NewModel()
		o := m.NewPlayer()
		l := m.NewRailLine(o)
		l.AutoExt = true

		n1 := m.NewRailNode(o, 0, 0)
		n2, e12 := n1.Extend(10, 0)
		head, _ := l.StartEdge(e12)

		_, e23 := n2.Extend(20, 0)

		TestCaseLineTasks{
			{"n1->n2", OnMoving, e12},
			{"n2->n3", OnMoving, e23},
			{"n3->n2", OnMoving, e23.Reverse},
			{"n2->n1", OnMoving, e12.Reverse},
			{"n1->n2", OnMoving, e12},
		}.Assert(t, head)

	})
}
