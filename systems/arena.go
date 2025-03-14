package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
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

func GetSpaceEntry(ecs *ecs.ECS) (spaceEntry *donburi.Entry) {
	spaceEntry, ok := tags.Space.First(ecs.World)
	if !ok {
		panic("No Space in ecs yet.")
	}
	return spaceEntry
}

func GetSpace(ecs *ecs.ECS) *resolv.Space {
	spaceEntry := GetSpaceEntry(ecs)
	return components.Space.Get(spaceEntry)
}

func GetArenaTileGrid(ecs *ecs.ECS) components.Grid {
	arenaEntry, _ := tags.Arena.First(ecs.World)
	// Grid is a slice, and slices internally use a pointer, to the
	// backing array. Thus we can pass it around directly
	// andthe values will contain the same address, thus chaning
	// the same underlying tilemap
	return components.TileGrid.Get(arenaEntry).Grid
}

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

// Remove wallEntry from Arena grid and from the
// collision space
// TODO: maybe dissolve walls slower
func BreakWall(wallEntry *donburi.Entry, ecs *ecs.ECS) {
	tile := components.Tile.Get(wallEntry)
	bbox := components.ConvexPolygonBBox.Get(wallEntry)
	x, y := tile.GridX, tile.GridY
	grid := GetArenaTileGrid(ecs)
	// remove from grid
	grid[x][y] = 0
	// remove Collision
	GetSpace(ecs).Remove(bbox)
}

func GetPlayer(ecs *ecs.ECS) *components.PlayerData {
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	return player
}

// The Arena is the 2d-grid where the player walks inside
// In here we update components that are represented by the
// 2d grid like breakable walls
func UpdateArena(ecs *ecs.ECS) {
	IncementTickCount(ecs)

	// if GetPlayer(ecs).Lives <= 0 {
	// 	fmt.Println("Game Over")
	// 	// TODO implement Game Over State / Screen
	// 	// Transition to Game Over State
	// }
	// Need some global timer to trigger the effect for a short time?
	//tileSpiralEffect(ecs)
	// handle breakable tiles/walls here
	//tileGrid := *GetArenaTileGrid(ecs)
	for entry := range tags.Tile.Iter(ecs.World) {
		bbox := components.ConvexPolygonBBox.Get(entry)
		if !bbox.Tags().Has(tags.TagBreakable) {
			continue
		}
		bbox.SelectTouchingCells(1).FilterShapes().
			ByTags(tags.TagExplosion).ForEach(
			func(shape resolv.IShape) bool {
				if bbox.Position().Equals(shape.Position()) {
					// breakable wall in explosion => dissolve
					BreakWall(entry, ecs)
					return false
				}
				return true
			},
		)
	}
	// for entry := range tags.Tile.Iter(ecs.World) {
	// 	bbox := components.ConvexPolygonBBox.Get(entry)
	// 	x := float32(bbox.Position().X - bbox.Bounds().Width()/2)
	// 	y := float32(bbox.Position().Y - bbox.Bounds().Width()/2)
	// 	vector.DrawFilledRect(screen, x, y, tileDiameter32, tileDiameter32, color.RGBA{0xff, 0, 0, uint8(entry.Id())}, false)
	// }

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
	backgroundTile := tileMap[0]
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
			// We draw the background tile first so that the tiles that have alpha
			// transparency show the background tile in their background
			screen.DrawImage(backgroundTile, op)
			screen.DrawImage(tileImage, op)
		}
	}

	// tileDiameter32 := float32(tileDiameter)
	// for entry := range tags.Tile.Iter(ecs.World) {
	// 	bbox := components.ConvexPolygonBBox.Get(entry)
	// 	x := float32(bbox.Position().X - bbox.Bounds().Width()/2)
	// 	y := float32(bbox.Position().Y - bbox.Bounds().Width()/2)
	// 	vector.DrawFilledRect(screen, x, y, tileDiameter32, tileDiameter32, color.RGBA{0xff, 0, 0, uint8(entry.Id())}, false)
	// }
}
