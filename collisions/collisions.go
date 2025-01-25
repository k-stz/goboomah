package collisions

import (
	"github.com/k-stz/goboomer/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

func AddCircleBBox(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		components.Space.Get(space).Add(GetCircleBBox(obj))
	}
}

func AddConvexPolygonBBox(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		components.Space.Get(space).Add(GetConvexPolygonBBox(obj))
	}
}

func GetCircleBBox(entry *donburi.Entry) *resolv.Circle {
	return components.CircleBBox.Get(entry)
}

func GetConvexPolygonBBox(entry *donburi.Entry) *resolv.ConvexPolygon {
	return components.ConvexPolygonBBox.Get(entry)
}
