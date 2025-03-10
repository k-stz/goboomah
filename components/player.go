package components

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type PlayerData struct {
	// represents in which direction and with what magnitude
	// the player wants to move
	// the collision/physics system will then calculate what
	// position this results into and will update the
	// players actual position (stored in another component)
	Direction resolv.Vector
	// Current movement Speed of the player, this is by how much
	// the player will move in a direction
	Movement resolv.Vector
	// how many bombs can be carried
	Bombs int
	// What Firepower placed bombs have
	Power int
	Lives int
	// When player loses a live respawn here
	RespawnPoint resolv.Vector
	// Whether Player is being damated, this is processed in an
	// FSM to see if the player needs to take damage or if
	// their are invincibility frames still
	Damaged bool
	// used for FSM implementation tracking Player damaged
	State State
	// How long to stay in a given state
	Duration TickCount
}

func NewPlayer(bombs, power int, lives int) *PlayerData {
	return &PlayerData{
		Bombs: bombs,
		Power: power,
		Lives: lives,
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
