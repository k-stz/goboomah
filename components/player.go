package components

import (
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
}

var Player = donburi.NewComponentType[PlayerData]()
