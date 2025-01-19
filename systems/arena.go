package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi/ecs"
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

// TODO draw all the tiles in here
// query all tiles and render them based on a tilemap, see stage.go
func DrawArena(ecs *ecs.ECS, screen *ebiten.Image) {
	// this is were I query the ecs to get all the data
	// to render the Stage (which is an "Arena")
	// The Arena itself gets created in factory.CreateArena
	for entry := range tags.Arena.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		o := components.TileGrid.Get(entry)
		drawColor := color.RGBA{180, 100, 0, 255}
		components.PrintGrid(o.Grid)
		fmt.Println("color", drawColor)
		//ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
	}
}

// func DrawFloatingPlatform(ecs *ecs.ECS, screen *ebiten.Image) {
// 	for e := range tags.FloatingPlatform.Iter(ecs.World) {
// 		o := dresolv.GetObject(e)
// 		drawColor := color.RGBA{180, 100, 0, 255}
// 		ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
// 	}
// }
