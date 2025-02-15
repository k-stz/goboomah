package factory

import (
	"fmt"

	"github.com/k-stz/goboomah/archtypes"
	"github.com/k-stz/goboomah/assets"
	"github.com/k-stz/goboomah/collisions"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/systems"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

func CreatePlayer(ecs *ecs.ECS, arenaEntry *donburi.Entry) *donburi.Entry {
	playerEntry := archtypes.Player.Spawn(ecs)
	// Set Player Data
	components.Player.Set(playerEntry, components.NewPlayer(10, 3))
	// Sprite
	playerSprite := components.Sprite.Get(playerEntry)
	playerSprite.Image = assets.Player

	w := float64(playerSprite.Image.Bounds().Dx())
	h := float64(playerSprite.Image.Bounds().Dy())

	tf := transform.Transform.Get(arenaEntry)
	dx := systems.GetWorldTileDiameter(ecs)
	x := tf.LocalPosition.X + (dx * 2.0)
	y := tf.LocalPosition.Y + (dx * 2.0)

	// BoundingCircle
	circleObj := resolv.NewCircle(x+w/2, y+h+w/2, w)
	circleObj = resolv.NewCircle(x, y, w/2)
	circleObj.Tags().Set(tags.TagPlayer)
	scale := 0.10
	components.ShapeCircle.Set(playerEntry, &components.ShapeCircleData{
		Circle:   circleObj,
		Radius:   w,
		Scale:    scale,
		Rotation: 0.0,
	})
	spaceEntry := systems.GetSpaceEntry(ecs)
	collisions.AddCircle(spaceEntry, playerEntry)
	fmt.Println("added player collision obj:", circleObj.Position())

	//dresolv.SetObject(platform, object)
	return playerEntry
}
