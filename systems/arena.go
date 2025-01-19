package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

// func UpdateFloatingPlatform(ecs *ecs.ECS) {
// 	for e := range tags.FloatingPlatform.Iter(ecs.World) {
// 		tw := components.Tween.Get(e)
// 		// Platform movement needs to be done first to make sure there's no space between the top and the player's bottom; otherwise, an alternative might
// 		// be to have the platform detect to see if the Player's resting on it, and if so, move the player up manually.
// 		y, _, seqDone := tw.Update(1.0 / 60.0)

// 		obj := dresolv.GetObject(e)
// 		obj.Y = float64(y)
// 		if seqDone {
// 			tw.Reset()
// 		}
// 	}
// }

// The Arena is the 2d-grid where the player walks inside
// For now we don't have any logic to update in here
func UpdateArena(ecs *ecs.ECS) {
	// the map itself is w
}

var mt = assets.Meadow_tile

// query all tiles and render them based on a tilemap, see stage.go
// Draw the stage which is referred to as an "Arena"
func DrawArena(ecs *ecs.ECS, screen *ebiten.Image) {
	// This is where we get the query the Arena Archetype
	// which contains the TileGrid, a 2d array of tileID
	// we will map the TileIDs to ebiten.Images to draw them
	//
	//
	// The Arena itself gets created in factory.CreateArena
	for entry := range tags.Arena.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		tg := components.TileGrid.Get(entry)
		tileMap := *components.TileMap.Get(entry)
		//components.PrintGrid(tg.Grid)
		tf := transform.Transform.Get(entry)
		dx := tg.TileDiameter
		var tileImage *ebiten.Image
		for x, row := range tg.Grid {
			for y, tileID := range row {
				// yes, looking up the image in a hash will kill locality
				// causing cache misses
				tileImage = tileMap[tileID]
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x)*dx, float64(y)*dx)
				op.GeoM.Translate(tf.LocalPosition.X, tf.LocalPosition.Y)
				screen.DrawImage(tileImage, op)
			}
		}
	}
}

// func DrawFloatingPlatform(ecs *ecs.ECS, screen *ebiten.Image) {
// 	for e := range tags.FloatingPlatform.Iter(ecs.World) {
// 		o := dresolv.GetObject(e)
// 		drawColor := color.RGBA{180, 100, 0, 255}
// 		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
// 	}
// }
