package factory

import (
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	playerEntry := archtypes.Player.Spawn(ecs)
	playerSprite := components.Sprite.Get(playerEntry)
	playerSprite.Image = assets.Player

	tf := transform.Transform.Get(playerEntry)
	tf.LocalPosition = math.Vec2{
		// TODO: this is a good idea, lets store screen dimensions
		// in a common game object (which in turn is in our ECS)
		// X: float64(game.Settings.ScreenWidth) * 0.75,
		//Y: cameraPos.Y + float64(game.Settings.ScreenHeight)*0.9,
		X: 50.0,
		Y: 200.0,
	}
	tf.LocalScale = math.Vec2{
		X: 0.5,
		Y: 0.5,
	}
	tf.LocalRotation = 0.0

	//dresolv.SetObject(platform, object)
	return playerEntry
}
