package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
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

// Cool spiral effect, just for asthetics.
// Use for stage transitions?
func tileSpiralEffect(ecs *ecs.ECS) {
	entry, ok := tags.Arena.First(ecs.World)
	if !ok {
		return
	}
	tf := transform.Transform.Get(entry)
	tf.LocalRotation = tf.LocalRotation + 0.01

}

// The Arena is the 2d-grid where the player walks inside
// For now we don't have any logic to update in here
func UpdateArena(ecs *ecs.ECS) {
	//tileSpiralEffect(ecs)

}

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

		var tileImage *ebiten.Image
		dx := tg.TileDiameter * tf.LocalScale.X
		var offsetX float64
		var offsetY float64
		for x, row := range tg.Grid {
			for y, tileID := range row {
				// yes, looking up the image in a hash will kill locality
				// causing cache misses
				tileImage = tileMap[tileID]
				offsetX = (float64(x) * dx) + tf.LocalPosition.X
				offsetY = float64(y)*dx + tf.LocalPosition.Y
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Rotate(tf.LocalRotation)
				op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
				op.GeoM.Translate(offsetX, offsetY)
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
