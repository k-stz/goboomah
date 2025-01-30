package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

func CreateBomb(position resolv.Vector, player *components.PlayerData, ecs *ecs.ECS) {
	bombEntry := archtypes.Bomb.Spawn(ecs)
	components.Bomb.Set(bombEntry, &components.BombData{
		Power:          player.Power,
		CountdownTicks: GetTickCount(ecs) + 200,
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
	position = SnapToGridPosition(position, dx)
	bbox := resolv.NewRectangle(position.X-dx/2, position.Y-dx/2, dx, dx)
	components.ConvexPolygonBBox.Set(bombEntry, bbox)
	fmt.Println("Bomb created", bombEntry.Id(), position)
}

// Attempt to place Bomb at Position by given player
func CanPlaceBombs(player *components.PlayerData) bool {
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
			ecs.World.Remove(entry.Entity())

		}

	}

}

// Snaps to the Center point of a grid
func SnapToGridPosition(pos resolv.Vector, tileDiameter float64) (newPosition resolv.Vector) {
	//pos.X = math.Round(pos.X/tileDiameter) * tileDiameter
	//pos.Y = math.Round(pos.Y/tileDiameter) * tileDiameter
	pos.X = math.Ceil(pos.X/tileDiameter) * tileDiameter
	pos.Y = math.Ceil(pos.Y/tileDiameter) * tileDiameter

	return pos
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
