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
	delegateField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)

	if !delegateField.IsValid() {
		return
	}

	if delegateField.IsNil() {
		var pid uint
		if parent := ch.Parent.Parent; parent != nil && parent.Data[oid] != nil {
			parentTarget := reflect.ValueOf(parent.Data[oid]).Elem().FieldByName(fieldName)
			pid = uint(parentTarget.Elem().FieldByName("ID").Uint())
		}
		delegate := reflect.New(delegateTypes[obj.B().T])
		delegate.Elem().FieldByName("DelegateNode").Set(reflect.ValueOf(ch.NewDelegateNode(obj, pid)))
		delegateField.Set(delegate)
	}
	delegateField.MethodByName("Add").Call([]reflect.Value{reflect.ValueOf(obj)})
}

func (ch *Chunk) addConnectable(obj Connectable) {
	chID := reflect.ValueOf(ch.ID)
	target := ch.M.RootCluster.FindChunk(obj.To(), ch.Parent.Scale)
	if target == nil {
		return
	}
	targetID := reflect.ValueOf(target.ID)
	outMapName := fmt.Sprintf("Out%ss", obj.B().Type().String())
	outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)

	if !outMap.IsValid() {
		return
	}

	if !outMap.MapIndex(targetID).IsValid() {
		delegate := reflect.New(delegateTypes[obj.B().T])
		delegate.Elem().FieldByName("DelegateEdge").Set(reflect.ValueOf(ch.NewDelegateEdge(obj, ch, target)))

		outMap.SetMapIndex(targetID, delegate)

		inMapName := fmt.Sprintf("In%ss", obj.B().Type().String())
		inMap := reflect.ValueOf(target).Elem().FieldByName(inMapName)
		inMap.SetMapIndex(chID, delegate)

		if _, ok := obj.(*RailEdge); ok {
			ch.setReverse(delegate.Interface().(*DelegateRailEdge), target)
		}
	}
	delegate := outMap.MapIndex(targetID)
	delegate.MethodByName("Add").Call([]reflect.Value{reflect.ValueOf(obj)})

}

func (ch *Chunk) setReverse(dre *DelegateRailEdge, target *Chunk) {
	if reverse, ok := target.OutRailEdges[ch.ID]; ok {
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
	delegateField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)

	if !delegateField.IsValid() {
		return
	}

	delegateField.MethodByName("Remove").Call([]reflect.Value{reflect.ValueOf(obj)})

	if delegateField.Elem().FieldByName("List").Len() == 0 {
		delegateField.Set(reflect.Zero(delegateField.Type()))
	}
}

func (ch *Chunk) removeConnectable(obj Connectable) {
	chID := reflect.ValueOf(ch.ID)
	target := ch.M.RootCluster.FindChunk(obj.To(), ch.Parent.Scale)
	if target == nil {
		return
	}
	targetID := reflect.ValueOf(target.ID)
	outMapName := fmt.Sprintf("Out%ss", obj.B().Type().String())
	outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)

	if !outMap.IsValid() {
		return
	}

	delegate := outMap.MapIndex(targetID)
	delegate.MethodByName("Remove").Call([]reflect.Value{reflect.ValueOf(obj)})

	if delegate.Elem().FieldByName("List").Len() == 0 {
		outMap.SetMapIndex(targetID, reflect.ValueOf(nil))
		inMapName := fmt.Sprintf("In%ss", obj.B().Type().String())
		inMap := reflect.ValueOf(target).Elem().FieldByName(inMapName)
		inMap.SetMapIndex(chID, reflect.ValueOf(nil))
	}
}

// Has returns whether specified Entity is deployed over Chunk or not.
func (ch *Chunk) Has(raw Entity) bool {
	id := reflect.ValueOf(raw.B().ID)
	switch obj := raw.(type) {
	case Localable:
		fieldName := obj.B().T.String()
		delegateField := reflect.ValueOf(ch).Elem().FieldByName(fieldName)
		if !delegateField.IsValid() || !delegateField.Elem().IsValid() {
			return false
		}
		return delegateField.Elem().FieldByName("List").MapIndex(id).IsValid()

	case Connectable:
		outMapName := fmt.Sprintf("Out%ss", obj.B().Type().String())
		outMap := reflect.ValueOf(ch).Elem().FieldByName(outMapName)
		if !outMap.IsValid() {
			return false
		}
		return outMap.MapIndex(id).IsValid()
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
	if rn := ch.RailNode; rn != nil {
		dm.RailNodes[rn.ID] = rn
	}
	for _, re := range ch.InRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.FromID] = re.From
	}
	for _, re := range ch.OutRailEdges {
		dm.RailEdges[re.ID] = re
		dm.RailNodes[re.ToID] = re.To
	}
}

// String represents status
func (ch *Chunk) String() string {
	return fmt.Sprintf("%s(%.1f:%d):u=%d,r=%v,c=%v,rn=%v,i=%d,o=%d:%v", ch.Type().Short(),
		ch.Parent.Scale, ch.ID, ch.OwnerID,
		ch.Residence, ch.Company, ch.RailNode,
		len(ch.InRailEdges), len(ch.OutRailEdges), ch.Point)
}
