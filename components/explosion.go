package components

import (
	"github.com/yohamta/donburi"
)

type ExplosionData struct {
	Power int // how many tiles the explosion extends
	// How long the explosion takes place
	CountdownTicks TickCount
	// you can call Draw() on it
}

var Explosion = donburi.NewComponentType[ExplosionData]()
