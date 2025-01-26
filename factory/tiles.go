package factory

import (
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

// create solid tiels from given arena entry
func CreateSolidTiles(ecs *ecs.ECS, arenaEntry *donburi.Entry) *donburi.Entry {
	tg := components.TileGrid.Get(arenaEntry)
	//diameter := tilegrid.TileDiameter

	tileMap := *components.TileMap.Get(arenaEntry)

	// walltile := tilemap[1]

	tf := transform.Transform.Get(arenaEntry)

	dx := tg.TileDiameter * tf.LocalScale.X
	tg.TileDiameter = dx
	for x, row := range tg.Grid {
		for y, tileID := range row {
			// 1 == is solid tiles
			if tileID == 1 {
				entry := archtypes.Tile.Spawn(ecs)
				components.Sprite.Get(entry).Image = tileMap[tileID]
				offsetX := (float64(x) * dx) + tf.LocalPosition.X
				offsetY := float64(y)*dx + tf.LocalPosition.Y
				bbox := resolv.NewRectangle(offsetX+dx/2, offsetY+dx/2, dx, dx)
				components.ConvexPolygonBBox.Set(entry, bbox)
			}
		}
	}
	//dresolv.SetObject(platform, object)
	return arenaEntry
}
