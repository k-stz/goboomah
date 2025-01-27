package components

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/math"
)

type PlayerData struct {
	// represents in which direction and with what magnitude
	// the player wants to move
	// the collision/physics system will then calculate what
	// position this results into and will update the
	// players actual position (stored in another component)
	Speed math.Vec2
	// how many bombs can be carried
	Bombs int
	// What Firepower placed bombs have
	Power int
}

func NewPlayer(bombs, power int) *PlayerData {
	return &PlayerData{
		Bombs: bombs,
		Power: power,
	}
}

// Used for Collision detection and also to superimpose the
// objects Image on top of it. That's why Scale, Radius, Rotation
// is added to derive the values when later drawing e.g. the
// players image on top of it.
type ShapeCircleData struct {
	Circle   *resolv.Circle
	Radius   float64 // original radius used for scaling
	Scale    float64 // Used to update Radius
	Rotation float64 // used to superimpose image inside
}

var Player = donburi.NewComponentType[PlayerData]()
var ShapeCircle = donburi.NewComponentType[ShapeCircleData]()
