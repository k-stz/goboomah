package components

import (
	"github.com/yohamta/donburi"
)

// Counts the ticks of the game
type TickCount int

type BombData struct {
	Power int
	// Game Tick count at which the bomb blows up
	CountdownTicks TickCount
	// whether the bomb shall blow up, for example
	// can be used to preempt the CountdownTicks in
	// chain-reaction explosion
	Detonate bool
}

var Tick = donburi.NewComponentType[TickCount]()
var Bomb = donburi.NewComponentType[BombData]()
