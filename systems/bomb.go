package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
)

func CreateBomb(position resolv.Vector, player *components.PlayerData, ecs *ecs.ECS) {
	bombEntry := archtypes.Bomb.Spawn(ecs)
	components.Bomb.Set(bombEntry, &components.BombData{
		Power:          player.Power,
		CountdownTicks: 500,
	})
	// Sprite
	components.Sprite.Set(bombEntry, &components.SpriteData{
		Image: assets.Bomb_tile,
	})
	// Shape
	bbox := resolv.NewRectangle(position.X, position.Y, 320, 320)
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

}

func DrawBomb(ecs *ecs.ECS, screen *ebiten.Image) {
	// TODO scale bombs
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

		var offsetX float64 = pos.X
		var offsetY float64 = pos.Y //- halfH + halfW

		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(4.0, 4.0)
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(bombSprite.Image, op)
	}
}
