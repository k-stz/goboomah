package factory

import (
	"github.com/k-stz/goboomah/archtypes"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
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
			if tileID == 0 {
				continue
			}
			entry := archtypes.Tile.Spawn(ecs)
			components.Sprite.Get(entry).Image = tileMap[tileID]
			offsetX := (float64(x) * dx) + tf.LocalPosition.X
			offsetY := float64(y)*dx + tf.LocalPosition.Y
			bbox := resolv.NewRectangle(offsetX+dx/2, offsetY+dx/2, dx, dx)
			switch tileID {
			// 1 == is solid tiles
			case 1:

				bbox.Tags().Set(tags.TagWall)
				// breakable wall
			case 2:
				bbox.Tags().Set(tags.TagWall | tags.TagBreakable)
			}
			components.ConvexPolygonBBox.Set(entry, bbox)
		}
	}
	//dresolv.SetObject(platform, object)
	return arenaEntry
}
