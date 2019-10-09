package entities

import (
	"fmt"
	"reflect"
)

// Chunk represents square area. Many Entities are deployed over Chunk.
type Chunk struct {
	Base
	Point

	Residence *DelegateResidence
	Company   *DelegateCompany
	RailNode  *DelegateRailNode

	Parent *Cluster

	InRailEdges  map[uint]*DelegateRailEdge
	OutRailEdges map[uint]*DelegateRailEdge
}

// NewChunk create Chunk on specified Cluster
func (m *Model) NewChunk(p *Cluster, o *Player) *Chunk {
	ch := &Chunk{
		Base:  m.NewBase(CHUNK, o),
		Point: p.Point,
	}
	ch.Init(m)
	ch.Resolve(p)
	m.Add(ch)
	return ch
}

// B returns base information of this elements.
func (ch *Chunk) B() *Base {
	return &ch.Base
}

// Init creates map.
func (ch *Chunk) Init(m *Model) {
	ch.Base.Init(CHUNK, m)
	ch.InRailEdges = make(map[uint]*DelegateRailEdge)
	ch.OutRailEdges = make(map[uint]*DelegateRailEdge)
}

// Add deploy Entity over Chunk
func (ch *Chunk) Add(raw Entity) {
	switch obj := raw.(type) {
	case Localable:
		ch.addLocalable(obj)
	case Connectable:
		ch.addConnectable(obj)
	}
}

func (ch *Chunk) addLocalable(obj Localable) {
	fieldName := obj.B().T.String()
	oid := obj.B().OwnerID
	nodeField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)

	if !nodeField.IsValid() {
		return
	}

	if nodeField.IsNil() {
		var pid uint
		if parent := ch.Parent.Parent; parent != nil && parent.Data[oid] != nil {
			parentTarget := reflect.ValueOf(parent.Data[oid]).Elem().FieldByName(fieldName)
			pid = uint(parentTarget.Elem().FieldByName("ID").Uint())
		}
		node := reflect.New(delegateTypes[obj.B().T])
		node.Elem().FieldByName("DelegateNode").Set(reflect.ValueOf(ch.NewDelegateNode(obj, pid)))
		nodeField.Set(node)
	}
	nodeField.MethodByName("Add").Call([]reflect.Value{reflect.ValueOf(obj)})
}

func (ch *Chunk) addConnectable(obj Connectable) {
	fromID := reflect.ValueOf(ch.ID)
	toCh := ch.M.RootCluster.FindChunk(obj.To(), ch.Parent.Scale)
	if toCh == nil {
		// ex. no Platform, Gate in Chunk referred by Step
		return
	}
	toID := reflect.ValueOf(toCh.ID)
	outMapName := fmt.Sprintf("Out%ss", obj.B().T.String())
	outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)

	if !outMap.IsValid() {
		// ex. no OutSteps in DelegateNode
		return
	}

	if !outMap.MapIndex(toID).IsValid() {
		nodeFieldName := connectTypes[obj.B().T].String()
		from := reflect.ValueOf(ch).Elem().FieldByName(nodeFieldName)
		to := reflect.ValueOf(toCh).Elem().FieldByName(nodeFieldName)

		edge := reflect.New(delegateTypes[obj.B().T])
		edge.Elem().FieldByName("DelegateEdge").Set(reflect.ValueOf(ch.NewDelegateEdge(
			obj, from.Interface().(delegateLocalable), to.Interface().(delegateLocalable))))

		outMap.SetMapIndex(toID, edge)

		inMapName := fmt.Sprintf("In%ss", obj.B().T.String())
		inMap := reflect.ValueOf(toCh).Elem().FieldByName(inMapName)
		inMap.SetMapIndex(fromID, edge)

		if _, ok := obj.(*RailEdge); ok {
			ch.setReverse(edge.Interface().(*DelegateRailEdge), toCh)
		}
	}
	edge := outMap.MapIndex(toID)
	edge.MethodByName("Add").Call([]reflect.Value{reflect.ValueOf(obj)})

}

func (ch *Chunk) setReverse(dre *DelegateRailEdge, toCh *Chunk) {
	if reverse, ok := toCh.OutRailEdges[ch.ID]; ok {
		dre.Reverse = reverse
		dre.ReverseID = reverse.ID
		reverse.Reverse = dre
		reverse.ReverseID = dre.ID
	}
}

// Remove undeploy Entity over Chunk
func (ch *Chunk) Remove(raw Entity) {
	switch obj := raw.(type) {
	case Localable:
		ch.removeLocalable(obj)
	case Connectable:
		ch.removeConnectable(obj)
	}
}

func (ch *Chunk) removeLocalable(obj Localable) {
	fieldName := obj.B().T.String()
	nodeField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)

	if !nodeField.IsValid() {
		return
	}

	nodeField.MethodByName("Remove").Call([]reflect.Value{reflect.ValueOf(obj)})

	if nodeField.Elem().FieldByName("List").Len() == 0 {
		nodeField.Set(reflect.Zero(nodeField.Type()))
	}
}

func (ch *Chunk) removeConnectable(obj Connectable) {
	fromID := reflect.ValueOf(ch.ID)
	toCh := ch.M.RootCluster.FindChunk(obj.To(), ch.Parent.Scale)
	if toCh == nil {
		// ex. no Platform, Gate in Chunk referred by Step
		return
	}
	toID := reflect.ValueOf(toCh.ID)
	outMapName := fmt.Sprintf("Out%ss", obj.B().T.String())
	outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)

	if !outMap.IsValid() {
		// ex. no OutSteps in DelegateNode
		return
	}

	delegate := outMap.MapIndex(toID)
	delegate.MethodByName("Remove").Call([]reflect.Value{reflect.ValueOf(obj)})

	if delegate.Elem().FieldByName("List").Len() == 0 {
		outMap.SetMapIndex(toID, reflect.ValueOf(nil))
		inMapName := fmt.Sprintf("In%ss", obj.B().T.String())
		inMap := reflect.ValueOf(toCh).Elem().FieldByName(inMapName)
		inMap.SetMapIndex(fromID, reflect.ValueOf(nil))
	}
}

// Has returns whether specified Entity is deployed over Chunk or not.
func (ch *Chunk) Has(raw Entity) bool {
	id := reflect.ValueOf(raw.B().ID)
	switch obj := raw.(type) {
	case Localable:
		fieldName := obj.B().T.String()
		nodeField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)
		if !nodeField.IsValid() || !nodeField.Elem().IsValid() {
			return false
		}
		return nodeField.Elem().FieldByName("List").MapIndex(id).IsValid()

	case Connectable:
		toCh := ch.M.RootCluster.FindChunk(obj.To(), ch.Parent.Scale)
		if toCh == nil {
			// ex. no Platform, Gate in Chunk referred by Step
			return false
		}
		toID := reflect.ValueOf(toCh.ID)

		outMapName := fmt.Sprintf("Out%ss", obj.B().T.String())
		outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)

		if !outMap.IsValid() {
			// ex. no OutSteps in DelegateNode
			return false
		}
		return outMap.MapIndex(toID).IsValid()
	}
	return false
}

// IsEmpty returns whether any Entity is deployed over Chunk or not.
func (ch *Chunk) IsEmpty() bool {
	return ch.RailNode == nil
}

// CheckDelete check remaining reference.
func (ch *Chunk) CheckDelete() error {
	return nil
}

// BeforeDelete remove reference of related entity
func (ch *Chunk) BeforeDelete() {
	ch.Parent.UnResolve(ch)
}

// Delete removes this entity with related ones.
func (ch *Chunk) Delete() {
	ch.M.Delete(ch)
}

// Resolve set reference
func (ch *Chunk) Resolve(args ...Entity) {
	for _, raw := range args {
		switch obj := raw.(type) {
		case *Cluster:
			ch.Parent = obj
		default:
			panic(fmt.Errorf("invalid type: %T %+v", obj, obj))
		}
	}
}

// Export set delegate Entity to DelegateMap
func (ch *Chunk) Export(dm *DelegateMap) {
	if r := ch.Residence; r != nil {
		dm.Residences[r.ID] = r
	}
	if c := ch.Company; c != nil {
		dm.Companies[c.ID] = c
	}
	if rn := ch.RailNode; rn != nil {
		dm.RailNodes[rn.ID] = rn
	}
	for _, re := range ch.InRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.FromID] = re.From.(*DelegateRailNode)
	}
	for _, re := range ch.OutRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.ToID] = re.To.(*DelegateRailNode)
	}
}

// String represents status
func (ch *Chunk) String() string {
	return fmt.Sprintf("%s(%.1f:%d):u=%d,r=%v,c=%v,rn=%v,i=%d,o=%d:%v", ch.T.Short(),
		ch.Parent.Scale, ch.ID, ch.OwnerID,
		ch.Residence, ch.Company, ch.RailNode,
		len(ch.InRailEdges), len(ch.OutRailEdges), ch.Point)
}
