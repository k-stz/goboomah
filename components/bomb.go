package components

import (
	"github.com/yohamta/donburi"
)

// Counts the ticks of the game
type TickCount int

type BombData struct {
	Power int
	// Time till bomb blows
	CountdownTicks int
}

var Tick = donburi.NewComponentType[TickCount]()
var Bomb = donburi.NewComponentType[BombData]()
