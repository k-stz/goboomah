package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

// Oh tweening is about inbetween animation
// this should be pretty useful!
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

func GetWorldTileDiameter(ecs *ecs.ECS) (tileDiameter float64) {
	// TODO for level extract Arena from GameScene object
	arenaEntry, _ := tags.Arena.First(ecs.World)
	return components.TileGrid.Get(arenaEntry).TileDiameter
}

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

func GetTickCount(ecs *ecs.ECS) components.TickCount {
	arenaEntry, _ := tags.Arena.First(ecs.World)
	return *components.Tick.Get(arenaEntry)
}

func IncementTickCount(ecs *ecs.ECS) {
	arenaEntry, _ := tags.Arena.First(ecs.World)
	*components.Tick.Get(arenaEntry)++
}

// The Arena is the 2d-grid where the player walks inside
// For now we don't have any logic to update in here
func UpdateArena(ecs *ecs.ECS) {
	IncementTickCount(ecs)

	// Need some global timer to trigger the effect for a short time?
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

	// Draw collision for all the
	arenaEntry, ok := tags.Arena.First(ecs.World)
	if !ok {
		panic("No arenaEntry")
	}
	// count game ticks, used for game logic like bomb fuse

	//o := dresolv.GetObject(e)
	tg := components.TileGrid.Get(arenaEntry)
	tileMap := *components.TileMap.Get(arenaEntry)
	//components.PrintGrid(tg.Grid)
	tf := transform.Transform.Get(arenaEntry)

	var tileImage *ebiten.Image
	tileDiameter := tg.TileDiameter
	var offsetX float64
	var offsetY float64
	for x, row := range tg.Grid {
		for y, tileID := range row {
			// yes, looking up the image in a hash will kill locality
			// causing cache misses
			tileImage = tileMap[tileID]
			offsetX = (float64(x) * tileDiameter) + tf.LocalPosition.X
			offsetY = float64(y)*tileDiameter + tf.LocalPosition.Y
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Rotate(tf.LocalRotation)
			op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
			//fmt.Println("tile x, y,", offsetX, offsetY)
			op.GeoM.Translate(offsetX, offsetY)
			screen.DrawImage(tileImage, op)
		}
	}

	tileDiameter32 := float32(tileDiameter)
	for entry := range tags.Tile.Iter(ecs.World) {
		// tf := transform.Transform.Get(entry)
		bbox := components.ConvexPolygonBBox.Get(entry)
		x := float32(bbox.Position().X - bbox.Bounds().Width()/2)
		y := float32(bbox.Position().Y - bbox.Bounds().Width()/2)
		//
		vector.DrawFilledRect(screen, x, y, tileDiameter32, tileDiameter32, color.RGBA{0xff, 0, 0, uint8(entry.Id())}, false)
		//x := tf.LocalPosition.X
		// y := tf.LocalPosition.Y
		// fmt.Println("x,y", x, y)
	}
}
