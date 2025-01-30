package factory

import (
	"fmt"

	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreatePlayer(ecs *ecs.ECS) *donburi.Entry {
	playerEntry := archtypes.Player.Spawn(ecs)
	// Set Player Data
	components.Player.Set(playerEntry, components.NewPlayer(10, 3))
	// Sprite
	playerSprite := components.Sprite.Get(playerEntry)
	playerSprite.Image = assets.Player

	w := float64(playerSprite.Image.Bounds().Dx())
	h := float64(playerSprite.Image.Bounds().Dy())

	x, y := 50.0, 170.0

	// BoundingCircle
	circleObj := resolv.NewCircle(x+w/2, y+h-w/2, w)
	circleObj.Tags().Set(tags.TagPlayer)
	scale := 0.20
	components.ShapeCircle.Set(playerEntry, &components.ShapeCircleData{
		Circle:   circleObj,
		Radius:   w,
		Scale:    scale,
		Rotation: 0.0,
	})
	fmt.Println("added player collision obj:", circleObj)

	//dresolv.SetObject(platform, object)
	return playerEntry
}
