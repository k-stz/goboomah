package tags

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var (
	Player    = donburi.NewTag().SetName("Player")
	Bomb      = donburi.NewTag().SetName("Bomb")
	Explosion = donburi.NewTag().SetName("Bomb")
	Arena     = donburi.NewTag().SetName("Arena")
	Tile      = donburi.NewTag().SetName("Tile")
	Space     = donburi.NewTag().SetName("Space")

	//Wall             = donburi.NewTag().SetName("Wall")
)

// Resolv tags

var (
	TagWall      = resolv.NewTag("Wall")
	TagBomb      = resolv.NewTag("Bomb")
	TagExplosion = resolv.NewTag("Explosion")
	TagPlayer    = resolv.NewTag("Player")
)
