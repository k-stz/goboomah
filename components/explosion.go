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

type ExplodableState struct {
	// Flag indicating whether object just suffered an explosion
	// but its impact hasn't been processed yet (this is where you
	// can decrement lives or start a death animation)
	ProcessedExplosion bool
	// Stay in Exploding state till we reach ExplodingTick
	ExplodingTick TickCount
	// Set to true while object sustained an explosion for till
	// ExplodingDuration Ticks is reached
	// Used to reacht to explosion: For example to hold still, play hit animation
	// Or to have invincibility frames
	Exploding bool
	// Target gametick Time till object despawns, can be used to show
	// death animation
	DespawnTick TickCount
	// When true objet and all its entries can be despawned/deleted
	Despawn bool
}

var Explodable = donburi.NewComponentType[ExplodableState]()
var Explosion = donburi.NewComponentType[ExplosionData]()
