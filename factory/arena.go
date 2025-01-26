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
	level := *components.LevelLa()
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
		X: 3.0,
		Y: 3.0,
	}
	tf.LocalRotation = 0

	//dresolv.SetObject(platform, object)
	return arenaEntry
}
