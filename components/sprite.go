package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// Will represent a particular tile/sprite that we will
// draw in tilemap cells or for the player. A map of pics,
// for easy implementation
type SpriteData struct {
	Image  *ebiten.Image
	Hidden bool
}

var Sprite = donburi.NewComponentType[SpriteData]()
