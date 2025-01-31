package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/collisions"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

func CreateDebugCircle(position resolv.Vector, radius float64, ecs *ecs.ECS) {
	debugCircleEntry := archtypes.DebugCircle.Spawn(ecs)
	// Sprite
	components.Sprite.Set(debugCircleEntry, &components.SpriteData{
		Image:  assets.Wall_tile,
		Hidden: true,
	})
	// Shape
	//dx := GetWorldTileDiameter(ecs)
	//position = SnapToGridPosition(position, dx)
	// below works, but it should work without the -dx/2
	//circleObj := resolv.NewCircle(position.X-dx/2, position.Y-dx/2, radius)
	circleObj := resolv.NewCircle(position.X, position.Y, radius)

	circleObj.Tags().Set(tags.TagDebug)
	components.ShapeCircle.Set(debugCircleEntry, &components.ShapeCircleData{
		Circle:   circleObj,
		Radius:   radius,
		Scale:    1.0,
		Rotation: 0.0,
	})
	spaceEntry := GetSpaceEntry(ecs)
	collisions.AddCircle(spaceEntry, debugCircleEntry)
	//fmt.Println("DebugCircle created", debugCircleEntry.Id(), position)
}

func UpdateDebugCircle(ecs *ecs.ECS) {
	if !ebiten.IsKeyPressed(ebiten.KeyT) {
		for entry := range tags.DebugCircle.Iter(ecs.World) {
			space := GetSpace(ecs)
			boundingCircleObject := components.ShapeCircle.Get(entry).Circle
			space.Remove(boundingCircleObject)
			ecs.World.Remove(entry.Entity())
		}
	}

	// currentGameTick := GetTickCount(ecs)
	// for entry := range tags.Explosion.Iter(ecs.World) {
	// 	explosion := components.Explosion.Get(entry)
	// 	if explosion.CountdownTicks <= currentGameTick {
	// 		fmt.Println("Blowing up!", entry.Entity())
	// 		ecs.World.Remove(entry.Entity())
	// 	}
	// 	// Handle collision logic here!
	// 	// Hurt players, items, walls? (movable walls)
	// }

}

func DrawDebugCircle(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.DebugCircle.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		debugSprite := components.Sprite.Get(entry)
		if debugSprite.Hidden {
			continue
		}

		halfW := float64(debugSprite.Image.Bounds().Dx() / 2)
		halfH := float64(debugSprite.Image.Bounds().Dy() / 2)

		circleObj := components.ShapeCircle.Get(entry).Circle
		pos := circleObj.Position()

		var offsetY float64 = pos.Y //- halfH + halfW
		var offsetX float64 = pos.X

		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		// Remove arena depending on scale
		arenaEntry, _ := tags.Arena.First(ecs.World)
		tf := transform.Transform.Get(arenaEntry)

		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(debugSprite.Image, op)
	}
}
