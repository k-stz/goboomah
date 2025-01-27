package components

import (
	"github.com/yohamta/donburi"
)

type BombData struct {
	Power int
	// Time till bomb blows
	CountdownTicks int
}

var Bomb = donburi.NewComponentType[BombData]()
