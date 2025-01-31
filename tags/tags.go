package tags

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var (
	Player      = donburi.NewTag().SetName("Player")
	Bomb        = donburi.NewTag().SetName("Bomb")
	Explosion   = donburi.NewTag().SetName("Bomb")
	Arena       = donburi.NewTag().SetName("Arena")
	Tile        = donburi.NewTag().SetName("Tile")
	Space       = donburi.NewTag().SetName("Space")
	DebugCircle = donburi.NewTag().SetName("DebugCircle")

	//Wall             = donburi.NewTag().SetName("Wall")
)

// Resolv tags

// var (
// 	TagWall      = resolv.NewTag("Wall")
// 	TagBomb      = resolv.NewTag("Bomb")
// 	TagExplosion = resolv.NewTag("Explosion")
// 	TagPlayer    = resolv.NewTag("Player")
// 	TagDebug     = resolv.NewTag("Debug")
// )

var (
	TagWall      = NewResolvTag("Wall")
	TagBomb      = NewResolvTag("Bomb")
	TagExplosion = NewResolvTag("Explosion")
	TagPlayer    = NewResolvTag("Player")
	TagDebug     = NewResolvTag("Debug")
)

var tagDirectory = map[resolv.Tags]string{}

// Used to get create resolv.Tag along with an internal
// Map so you can call .String() on it with shorter result
func NewResolvTag(tagName string) resolv.Tags {
	tag := resolv.NewTag(tagName)
	tagDirectory[tag] = tagName
	return tag
}

func ToString(tag resolv.Tags) string {
	// TODO add "comma, ok" idiom check, return 
	// empty string then
	return tagDirectory[tag]
}
