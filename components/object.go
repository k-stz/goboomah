package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
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

// Circle Bounding Box. Becuase can't use a generic container
// like resolv.IShape becuase the "NewComponentType" creates a
// pointer of it and at this point I can't put things inside it
var CircleBBox = donburi.NewComponentType[resolv.Circle]()
var ConvexPolygonBBox = donburi.NewComponentType[resolv.ConvexPolygon]()

// doesn't work...
//var Object = donburi.NewComponentType[resolv.IShape]()

func SetCircleBBox(circle *resolv.Circle, tf *transform.TransformData, sprite *ebiten.Image) {
	halfW := float64(sprite.Bounds().Dx() / 2)
	h := float64(sprite.Bounds().Dy())

	x := tf.LocalPosition.X
	y := tf.LocalPosition.Y

	scaleX := tf.LocalScale.X
	circle.SetRadius(halfW * scaleX)

	newX := x + halfW //*scaleX
	newY := y + h - halfW
	circle.SetPosition(newX, newY)
}

func CircleBottomLeftPos(c *resolv.Circle) resolv.Vector {
	x := c.Position().X - c.Radius()
	y := c.Position().Y - c.Radius()
	return resolv.NewVector(x, y)
}
