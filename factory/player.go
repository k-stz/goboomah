package factory

import (
	"fmt"

	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	playerEntry := archtypes.Player.Spawn(ecs)
	playerSprite := components.Sprite.Get(playerEntry)
	playerSprite.Image = assets.Player

	w := float64(playerSprite.Image.Bounds().Dx())
	h := float64(playerSprite.Image.Bounds().Dy())

	x, y := -20.0, 200.0

	circleObj := resolv.NewCircle(x+w/2, y+h-w/2, w)
	//components.SetCircleBBox(circleObj, tf, playerSprite.Image)
	scale := 1.0
	components.ShapeCircle.Set(playerEntry, &components.ShapeCircleData{
		Circle:   circleObj,
		Radius:   w,
		Scale:    scale,
		Rotation: 0.0,
	})

	// TODO deltete all transform.Transform references once this is done
	tf := transform.Transform.Get(playerEntry)
	tf.LocalPosition = math.Vec2{
		// TODO: this is a good idea, lets store screen dimensions
		// in a common game object (which in turn is in our ECS)
		// X: float64(game.Settings.ScreenWidth) * 0.75,
		//Y: cameraPos.Y + float64(game.Settings.ScreenHeight)*0.9,
		X: x,
		Y: y,
	}
	tf.LocalScale = math.Vec2{
		X: scale,
		Y: scale,
	}
	tf.LocalRotation = 0.0

	fmt.Println("added player collision obj:", circleObj)

	//dresolv.SetObject(platform, object)
	return playerEntry
}
