package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/archtypes"
	"github.com/k-stz/goboomah/assets"
	"github.com/k-stz/goboomah/collisions"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

func CreateBomb(position resolv.Vector, player *components.PlayerData, ecs *ecs.ECS) {
	bombEntry := archtypes.Bomb.Spawn(ecs)
	components.Bomb.Set(bombEntry, &components.BombData{
		Power:          player.Power,
		CountdownTicks: GetTickCount(ecs) + 75,
		Detonate:       false,
	})
	// Sprite
	components.Sprite.Set(bombEntry, &components.SpriteData{
		//Image: assets.Bomb_tile,
		Image: assets.Bomb_tile,
	})
	// Shape
	dx := GetWorldTileDiameter(ecs)
	fmt.Println("before snap pos:", position)
	position = SnapToGridTileCenter(position, dx)
	bbox := resolv.NewRectangle(position.X, position.Y, dx, dx)
	bbox.Tags().Set(tags.TagBomb)
	components.ConvexPolygonBBox.Set(bombEntry, bbox)
	// Add Shape to space
	collisions.AddConvexPolygonBBox(GetSpaceEntry(ecs), bombEntry)
	fmt.Println("Bomb created", bombEntry.Id(), position)
}

// Attempt to place Bomb at Position by given player
func CanPlaceBombs(checkPosition resolv.Vector, ecs *ecs.ECS) bool {
	dx := GetWorldTileDiameter(ecs)
	tileShapeTags, _ := CheckTile(checkPosition, dx/2, false, ecs)
	if tileShapeTags.Has(tags.TagBomb | tags.TagWall) {
		return false
	}

	return true
	// implement logic later
	//return player.Bombs > 0
}

// Update bomb ticks
func UpdateBomb(ecs *ecs.ECS) {
	currentGameTick := GetTickCount(ecs)
	for entry := range tags.Bomb.Iter(ecs.World) {
		bomb := components.Bomb.Get(entry)
		bombPosition := components.ConvexPolygonBBox.Get(entry).Position()
		if bomb.CountdownTicks <= currentGameTick {
			// We set the bomb to exploding that's how we
			// can later add other logic to make a bomb explode sooner
			// blow up bomb
			bomb.Detonate = true

		}
		if bomb.Detonate {
			fmt.Println("Blowing up!", entry.Entity())
			CreateExplosion(bombPosition, bomb.Power, ecs)
			RemoveBomb(entry, ecs)
		}
	}
}

// Remove Bomb from ecs and its object from the Collisoin spacce
func RemoveBomb(bombEntry *donburi.Entry, ecs *ecs.ECS) {
	space := GetSpace(ecs)
	bbox := components.ConvexPolygonBBox.Get(bombEntry)
	space.Remove(bbox)
	ecs.World.Remove(bombEntry.Entity())
}

// "Snaps into place" the position vector to the grid's tile
// where the grid's tile has a diameter of "tileDiameter".
// In whatever tile "pos"ition falls, the top-left corner will be returned.
// if the grid starts and 0,0 and each tile/cell of the grid 40.0x40.0
// then for the postion (20,30) it will return (0,0)
// for the postion (50, 70) it will return (40,40)
func SnapToGridTileTopLeft(position resolv.Vector, tileDiameter float64) (newPosition resolv.Vector) {
	position.X = math.Floor(position.X/tileDiameter) * tileDiameter
	position.Y = math.Floor(position.Y/tileDiameter) * tileDiameter

	return position
}

// Like SnapToGridTileTopLeft but will snap in the center of a Tile.
// if the grid starts at 0,0 and each tile/cell of the grid 40.0x40.0
// then for the postion (20,30) it will return (20,20)
// for the postion (50, 70) it will return (60,60)
func SnapToGridTileCenter(position resolv.Vector, tileDiameter float64) (newPosition resolv.Vector) {
	dx := tileDiameter
	position.X = (math.Floor(position.X/dx) * dx) + dx/2
	position.Y = (math.Floor(position.Y/dx) * dx) + dx/2

	return position
}

func DrawBomb(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Bomb.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		bombSprite := components.Sprite.Get(entry)

		halfW := float64(bombSprite.Image.Bounds().Dx() / 2)
		halfH := float64(bombSprite.Image.Bounds().Dy() / 2)

		bbox := components.ConvexPolygonBBox.Get(entry)
		pos := bbox.Position()
		rotation := bbox.Rotation()
		// rad := halfW
		// diameter := max(halfW, halfH)
		// // diameter * x = radius
		// scale := rad / diameter

		var offsetY float64 = pos.Y //- halfH + halfW
		var offsetX float64 = pos.X

		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		// Remove arena depending on scale
		arenaEntry, _ := tags.Arena.First(ecs.World)
		tf := transform.Transform.Get(arenaEntry)

		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(rotation)
		op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(bombSprite.Image, op)
	}
}
