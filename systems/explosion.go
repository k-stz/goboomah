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
	"github.com/yohamta/donburi/features/transform"
)

func CreateExplosion(position resolv.Vector, reach int, ecs *ecs.ECS) {
	explosionEntry := archtypes.Explosion.Spawn(ecs)
	components.Explosion.Set(explosionEntry, &components.ExplosionData{
		Power:          reach,
		CountdownTicks: GetTickCount(ecs) + 100,
	})
	// Sprite
	components.Sprite.Set(explosionEntry, &components.SpriteData{
		Image: assets.Wall_tile,
	})
	// Shape
	dx := GetWorldTileDiameter(ecs)
	position = SnapToGridPosition(position, dx)
	bbox := resolv.NewRectangle(position.X-dx/2, position.Y-dx/2, dx, dx)
	components.ConvexPolygonBBox.Set(explosionEntry, bbox)
	fmt.Println("Bomb created", explosionEntry.Id(), position)
}

func UpdateExplosion(ecs *ecs.ECS) {
	currentGameTick := GetTickCount(ecs)
	for entry := range tags.Explosion.Iter(ecs.World) {
		explosion := components.Explosion.Get(entry)
		if explosion.CountdownTicks <= currentGameTick {
			fmt.Println("Blowing up!", entry.Entity())
			ecs.World.Remove(entry.Entity())
		}
		// Handle collision logic here!
		// Hurt players, items, walls? (movable walls)
	}

}

func DrawExplosion(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Explosion.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		explosionSprite := components.Sprite.Get(entry)

		halfW := float64(explosionSprite.Image.Bounds().Dx() / 2)
		halfH := float64(explosionSprite.Image.Bounds().Dy() / 2)

		bbox := components.ConvexPolygonBBox.Get(entry)
		pos := bbox.Position()
		rotation := bbox.Rotation()

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
		screen.DrawImage(explosionSprite.Image, op)
	}
}
