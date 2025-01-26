package collisions

import (
	"github.com/k-stz/goboomer/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

func AddCircle(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		components.Space.Get(space).Add(GetCircleShape(obj))
	}
}

func AddConvexPolygonBBox(space *donburi.Entry, objects ...*donburi.Entry) {
	for _, obj := range objects {
		components.Space.Get(space).Add(GetConvexPolygonBBox(obj))
	}
}

func GetCircleShape(entry *donburi.Entry) *resolv.Circle {
	return components.ShapeCircle.Get(entry).Circle
}

func GetConvexPolygonBBox(entry *donburi.Entry) *resolv.ConvexPolygon {
	return components.ConvexPolygonBBox.Get(entry)
}
