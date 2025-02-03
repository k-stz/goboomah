package components

import (
	"github.com/yohamta/donburi"
	"github.com/yohamta/ganim8/v2"
)

type ExplosionData struct {
	Power int // how many tiles the explosion extends
	// How long the explosion takes place
	CountdownTicks TickCount
	explosion      *ganim8.Animation
}

var Explosion = donburi.NewComponentType[ExplosionData]()
