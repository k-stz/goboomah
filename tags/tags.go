package tags

import (
	"strconv"

	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

var (
	CircleCollision = donburi.NewTag().SetName("CircleCollsion")
	Player          = donburi.NewTag().SetName("Player")
	Enemy           = donburi.NewTag().SetName("Enemey")
	Bomb            = donburi.NewTag().SetName("Bomb")
	Explosion       = donburi.NewTag().SetName("Bomb")
	Arena           = donburi.NewTag().SetName("Arena")
	Tile            = donburi.NewTag().SetName("Tile")
	Space           = donburi.NewTag().SetName("Space")
	DebugCircle     = donburi.NewTag().SetName("DebugCircle")

	//Wall             = donburi.NewTag().SetName("Wall")
)

// Resolv tags

var (
	TagWall = NewResolvTag("Wall")
	// enemies shall hurt on touch, while players shall walk
	// through eachother, so different collision tags used
	TagEnemy     = NewResolvTag("TagEnemy")
	TagBreakable = NewResolvTag("Breakable")
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

// Returns string representation of tag, if tag doesn't exists
// its numeric uint representation is returned
// Adapted from resolv's lib "func (t Tags) String() string" implementation
// to make it less verbose to build debugging on top
func ToString(tag resolv.Tags) string {
	result := ""

	tagIndex := 0

	for i := 0; i < 64; i++ {
		possibleTag := resolv.Tags(1 << i)
		if tag.Has(possibleTag) {
			if tagIndex > 0 {
				result += "|"
			}

			value, ok := tagDirectory[possibleTag]

			if !ok {
				value = strconv.Itoa(int(possibleTag))
			}

			result += value
			tagIndex++
		}
	}

	return result
}
