package tags

import "github.com/yohamta/donburi"

var (
	Player = donburi.NewTag().SetName("Player")
	Bomb   = donburi.NewTag().SetName("Bomb")
	Arena  = donburi.NewTag().SetName("Arena")
	Tile   = donburi.NewTag().SetName("Tile")
	Space  = donburi.NewTag().SetName("Space")

	//Wall             = donburi.NewTag().SetName("Wall")
)
