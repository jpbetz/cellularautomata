// automatically generated by the FlatBuffers compiler, do not modify

package region

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type GuardUnit struct {
	_tab flatbuffers.Table
}

func GetRootAsGuardUnit(buf []byte, offset flatbuffers.UOffsetT) *GuardUnit {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &GuardUnit{}
	x.Init(buf, n+offset)
	return x
}

func (rcv *GuardUnit) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *GuardUnit) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *GuardUnit) Position(obj *Position) *Position {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		x := o + rcv._tab.Pos
		if obj == nil {
			obj = new(Position)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *GuardUnit) Waypoints(obj *Position, j int) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		x := rcv._tab.Vector(o)
		x += flatbuffers.UOffsetT(j) * 8
		obj.Init(rcv._tab.Bytes, x)
		return true
	}
	return false
}

func (rcv *GuardUnit) WaypointsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func GuardUnitStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func GuardUnitAddPosition(builder *flatbuffers.Builder, position flatbuffers.UOffsetT) {
	builder.PrependStructSlot(0, flatbuffers.UOffsetT(position), 0)
}
func GuardUnitAddWaypoints(builder *flatbuffers.Builder, waypoints flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(waypoints), 0)
}
func GuardUnitStartWaypointsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 4)
}
func GuardUnitEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
