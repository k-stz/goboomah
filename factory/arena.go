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

type TileGrid [][]int

func CreateArena(ecs *ecs.ECS) *donburi.Entry {
	//level1 := components.Level1()
	arenaEntry := archtypes.Arena.Spawn(ecs)
	// TODO create tile objects based on TileGrid!
	level := *components.Level1()
	components.TileGrid.SetValue(arenaEntry, level)

	// Creat Tile Mapping, used for dynamic render loop
	level1TileMap := components.TileMapData{
		0: assets.Meadow_tile,
		1: assets.Wall_tile,
	}
	components.TileMap.SetValue(arenaEntry, level1TileMap)

	tf := transform.Transform.Get(arenaEntry)
	tf.LocalPosition = math.Vec2{
		// TODO: this is a good idea, lets store screen dimensions
		// in a common game object (which in turn is in our ECS)
		// X: float64(game.Settings.ScreenWidth) * 0.75,
		//Y: cameraPos.Y + float64(game.Settings.ScreenHeight)*0.9,
		X: 100.0,
		Y: 100.0,
	}
	tf.LocalScale = math.Vec2{
		X: 4.0,
		Y: 4.0,
	}
	tf.LocalRotation = 90.0

	//dresolv.SetObject(platform, object)
	return arenaEntry
}

// func CreateFloatingPlatform(ecs *ecs.ECS, object *resolv.Object) *donburi.Entry {
// 	platform := archetypes.FloatingPlatform.Spawn(ecs)
// 	dresolv.SetObject(platform, object)

// 	// The floating platform moves using a *gween.Sequence sequence of tweens, moving it back and forth.
// 	tw := gween.NewSequence()
// 	obj := components.Object.Get(platform)
// 	tw.Add(
// 		gween.New(float32(obj.Y), float32(obj.Y-128), 2, ease.Linear),
// 		gween.New(float32(obj.Y-128), float32(obj.Y), 2, ease.Linear),
// 	)
// 	components.Tween.Set(platform, tw)

// 	return platform
// }
