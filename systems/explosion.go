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
	position = SnapToGridTileCenter(position, dx)
	bbox := resolv.NewRectangle(position.X, position.Y, dx, dx)
	bbox.Tags().Set(tags.TagExplosion)
	components.ConvexPolygonBBox.Set(explosionEntry, bbox)
	fmt.Println("Bomb created", explosionEntry.Id(), position)
}

type TileContent struct {
	CenterPosition      resolv.Vector // the snapped center position of tile
	CollisionObjectTags resolv.Tags   // tags of objects found in the tile
	IsEmpty             bool
}

type Direction int

const (
	Up Direction = iota + 1
	Down
	Left
	Right
)

// Returns resolv.Vector pointing in direction by "length"
// in Windowcoordinates (negative y is up).
// Example: GetDirectionVector(Up, 40.0) => (0, -40.0)
func GetDirectionVector(direction Direction, length float64) resolv.Vector {
	vec := resolv.NewVectorZero()
	switch direction {
	case Up:
		vec.Y = -length
	case Down:
		vec.Y = length
	case Right:
		vec.X = length
	case Left:
		vec.X = -length
	}
	return vec
}

// Input:
// Checks tileCount number of Tiles (inclusive) "fromPos" in "direction"
// tileDiameter. On "fromPos" SnapToTileGridCenter is applied for tile alignment
// Inside the Tiles it is tested if collision objects with the given resolv.Tags are present
// debugMode: Whether to create debugCircle to visualize the checks for debugging
//
// Returns: For each tile if any tags where found inside and their position
func CheckTilesInDirection(fromPos resolv.Vector, direction Direction, tileCount int, tileDiameter float64,
	forTags resolv.Tags, debugMode bool, ecs *ecs.ECS) (tileContents []TileContent) {
	dx := tileDiameter
	pos := SnapToGridTileCenter(fromPos, dx)
	checkTiles := tileCount
	tileContents = []TileContent{}
	dirVector := GetDirectionVector(direction, dx)
	for i := range checkTiles {
		//offsetY := -(dx + (float64(i) * dx))
		checkPos := pos.Add(dirVector.Scale(float64(i + 1)))
		tileShapeTags, _ := CheckTile(checkPos, dx/2, debugMode, ecs)

		tileResult := TileContent{
			CenterPosition:      checkPos,
			CollisionObjectTags: 0,
			IsEmpty:             true,
		}
		for _, shapeTag := range tileShapeTags {
			if shapeTag.Has(forTags) {
				tileResult.IsEmpty = false
				tileResult.CollisionObjectTags = shapeTag
				break
			}
		}
		tileContents = append(tileContents, tileResult)
	}
	return tileContents
}

// Input:
// debugMode: whether to create a debugCircle Object to draw the circle used for checking
// for debugging
// Returns:
// tileShapeTags shapes that were found at the checkPosition in the radius
// isEmpty: indicates whether any tags were found
func CheckTile(checkPosition resolv.Vector, radius float64, debugMode bool, ecs *ecs.ECS) (tileShapeTags []resolv.Tags, isEmpty bool) {
	spaceEntry, _ := tags.Space.First(ecs.World)
	space := components.Space.Get(spaceEntry)

	if debugMode {
		CreateDebugCircle(checkPosition, radius, ecs)
	}

	// tc is a tile checker, which is a circle bounding box object
	// used to scan over tiles for intersection to if it intersects in a particular
	// tile with objects of interest
	tc := resolv.NewCircle(checkPosition.X, checkPosition.Y, radius)
	space.Add(tc)
	defer space.Remove(tc) // remove circle scanner

	//intersectionFound :=
	tc.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: tc.SelectTouchingCells(1).FilterShapes().ByTags(tags.TagWall),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			// lets report what we touched with rather
			//playerShape.Circle.MoveVec(set.MTV)
			// set.OtherShape = The other shape involved in the contact.
			// iscontainedby is what we need here!
			insideWall := checkPosition.IsInside(set.OtherShape)
			if insideWall {
				tileShapeTags = append(tileShapeTags, *set.OtherShape.Tags())
				return false // stop testing for further intersection
			}
			//set.OtherShape.(*resolv.ConvexPolygon).IsContainedBy(set.OtherShape)
			//fmt.Println("COLLISION tag", set.OtherShape)
			return true
		},
	})
	isEmpty = true
	if len(tileShapeTags) > 0 {
		isEmpty = false
	}

	return tileShapeTags, isEmpty
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
