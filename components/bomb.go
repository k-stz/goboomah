package components

import (
	"github.com/yohamta/donburi"
)

// Counts the ticks of the game
type TickCount int

type BombData struct {
	Power int
	// Game Tick count at which the bomb blows up
	// Detonate is set to true
	CountdownTicks TickCount
	// After "Explode" is set, short delay before bomb goes
	// off for real
	ExplosionDelayTicks TickCount
	// whether the bomb shall blow up, for example
	// can be used to preempt the CountdownTicks in
	// chain-reaction explosion
	Explode bool
	// when bomb shall actually despawn, after which
	// it can be replaced by an explosion
	Despawn bool
}

var Tick = donburi.NewComponentType[TickCount]()
var Bomb = donburi.NewComponentType[BombData]()
