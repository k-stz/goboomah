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
	bbox.Tags().Set(tags.TagExplosion)
	components.ConvexPolygonBBox.Set(explosionEntry, bbox)
	fmt.Println("Bomb created", explosionEntry.Id(), position)
}

// func TileChecker(from resolv.Vector, to resolv.Vector, ecs *ecs.ECS) {
// Returns all shapes intersection roughly with the tile
func CheckTile(checkPosition resolv.Vector, ecs *ecs.ECS) []resolv.Tags {
	//fmt.Printf("CheckTile at %v\n", checkPosition)
	spaceEntry, _ := tags.Space.First(ecs.World)
	space := components.Space.Get(spaceEntry)

	dx := GetWorldTileDiameter(ecs)
	// tc is a tile checker, which is a circle bounding box object
	// used to scan over tiles for intersection to if it intersects in a particular
	// tile with objects of interest
	// For now dx/2 in diameter to not accidentally touch neighboring
	//
	tc := resolv.NewCircle(checkPosition.X, checkPosition.Y, dx/2)
	space.Add(tc)
	defer space.Remove(tc) // remove circle scanner

	// TODO check if resolv.LineTest is better
	var shapeTags []resolv.Tags
	//intersectionFound :=
	tc.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: tc.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			// lets report what we touched with rather
			//playerShape.Circle.MoveVec(set.MTV)
			// set.OtherShape = The other shape involved in the contact.
			shapeTags = append(shapeTags, *set.OtherShape.Tags())
			//fmt.Println("COLLISION tag", set.OtherShape)
			return true
		},
	})
	// if intersectionFound {
	// 	fmt.Println("collision found at:", checkPosition)
	// 	//fmt.Println("shapes:", shapes)
	// 	for i, s := range shapes {
	// 		fmt.Printf("%d tag: %v\n", i, s.Tags())
	// 	}
	// } else {
	// 	fmt.Println("nothing")
	// }

	return shapeTags
	//space.ForEachShape()

	// test a selectino of shapes against a line
	//resolv.LineTest(lts)
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
