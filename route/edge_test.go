package route

import "testing"

func TestEdge(t *testing.T) {
	t.Run("Export", func(t *testing.T) {
		org1 := &Node{}
		org2 := &Node{}
		orgE := &Edge{}
		orgE.FromNode, orgE.ToNode = org1, org2

		new1 := &Node{}
		new1.Out = []*Edge{}
		new2 := &Node{}
		new2.In = []*Edge{}
		newE := orgE.Export(new1, new2)

		if orgE == newE {
			t.Errorf("newE wants not %v, but does", orgE)
		}

		if newE.FromNode != new1 {
			t.Errorf("newE.FromNode wants %v, but not", new1)
		}
		if newE.ToNode != new2 {
			t.Errorf("newE.ToNode wants %v, but not", new2)
		}

		if new1.Out[0] != newE {
			t.Errorf("new1.Out[0] wants to %v, but not", newE)
		}
		if new1.Out[0].FromNode != new1 {
			t.Errorf("new1.Out[0].FromNode wants to %v, but not", new1)
		}
		if new1.Out[0].ToNode != new2 {
			t.Errorf("new1.Out[0].ToNode wants to %v, but not", new2)
		}

		if new2.In[0] != newE {
			t.Errorf("new2.In[0] wants to %v, but not", newE)
		}
		if new2.In[0].FromNode != new1 {
			t.Errorf("new2.In[0].FromNode wants to %v, but not", new1)
		}
		if new2.In[0].ToNode != new2 {
			t.Errorf("new2.In[0].ToNode wants to %v, but not", new2)
		}
	})
}
