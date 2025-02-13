package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type RectangleData struct {
	X, Y, W, H float64
}

func NewRectangle(x, y, w, h float64) *RectangleData {
	return &RectangleData{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

// indicates whether something something is walkable (tree vs bush)
// not used
type CollidableData struct {
	IsSolid bool
}

// I can give it a default value in the parenthesis here..
var Rectangle = donburi.NewComponentType[RectangleData]()
var Image = donburi.NewComponentType[ebiten.Image]()
var Collidable = donburi.NewComponentType[CollidableData]()

// used for collision detection/response
// IShape is an interface that can fit any shape
// like ConvexRectangle and Circle which should
// be the only shapes we need

var ConvexPolygonBBox = donburi.NewComponentType[resolv.ConvexPolygon]()
