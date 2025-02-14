package factory

import (
	"github.com/hajimehoshi/ebiten/v2"
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

// Spawn enemy tiles from given arena entry
func CreateEnemies(ecs *ecs.ECS, arenaEntry *donburi.Entry) *donburi.Entry {
	tg := components.TileGrid.Get(arenaEntry)
	//diameter := tilegrid.TileDiameter

	// walltile := tilemap[1]

	tf := transform.Transform.Get(arenaEntry)

	dx := systems.GetWorldTileDiameter(ecs)
	for x, row := range tg.Grid {
		for y, tileID := range row {
			// lets say tileID bigger than 5 will be enemies
			if tileID <= 5 {
				continue
			}
			entry := archtypes.Enemy.Spawn(ecs)
			// tile grid data
			// tileData := components.Tile.Get(entry)
			// tileData.GridX, tileData.GridY = x, y
			// Sprite
			var enemySprite *ebiten.Image
			switch tileID {
			case 9:
				enemySprite = assets.Blob_npc
			}
			components.Sprite.Get(entry).Image = enemySprite
			// Collisions
			offsetX := (float64(x) * dx) + tf.LocalPosition.X
			offsetY := float64(y)*dx + tf.LocalPosition.Y

			w := float64(enemySprite.Bounds().Dx())
			//h := float64(enemySprite.Bounds().Dy())

			radius := (dx / 2)

			//boudingCircle := resolv.NewCircle(offsetX+w/2, offsetY+h-w/2, w)
			boudingCircle := resolv.NewCircle(offsetX+dx/2, offsetY+dx/2, radius)
			//bbox := resolv.NewRectangle(offsetX+dx/2, offsetY+dx/2, dx, dx)
			switch tileID {
			// 1 == is solid tiles
			//  consolidate with CreateTile func?
			case 9:
				boudingCircle.Tags().Set(tags.TagEnemy | tags.TagDebug)
			}
			components.ShapeCircle.Set(entry, &components.ShapeCircleData{
				Circle: boudingCircle,
				Radius: radius, // dx / 2
				// scale image to fit bouding circle in width
				Scale:    dx / w,
				Rotation: 0.0,
			})
			spaceEntry := systems.GetSpaceEntry(ecs)
			collisions.AddCircle(spaceEntry, entry)
		}
	}
	//dresolv.SetObject(platform, object)
	return arenaEntry
}
