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
	"github.com/yohamta/ganim8/v2"
)

// Returns all CenterPositions in TileContents in order till non empty one is
// found
// Used for example when tryig to figure out how many tiles an explosion shall
// cover till it stops at a wall (=non-empty; thus have an non-empty TileContents)
// Example: When tileContents [val1 val2 empty val4] then it returns
// [val1.CenterPosition val2.CenterPosition]
//
// TODO: Could also be improved to stop at a specific collision tag!
func TakeUntilNonEmpty(tileContnets []TileContent) []resolv.Vector {
	var positions []resolv.Vector
	fmt.Println("TakeUntilEmpty Loop")
	for _, v := range tileContnets {
		fmt.Println("positoin", positions)
		if !v.IsEmpty {
			return positions
		}
		positions = append(positions, v.CenterPosition)
	}
	return positions
}

// Calculates all position were an explosion shall be spawned, using "atPosition"
// as the origin of the explosion (here the bomb was basically placed)
// Inputs:
// atPostion: origin of the postions, should be at the center of a tile
// reach: the explosion reach of a bomb
// Returns:
// Slice containing 4 slices, for each cardinal postions
// such that we can interpolate different animations frames for the explosions
// spawned for each offset away from the bomb!
func GetExplosionPositions(atPosition resolv.Vector, reach int, ecs *ecs.ECS) (explostionDirectionPositions [][]resolv.Vector) {
	dx := GetWorldTileDiameter(ecs)
	var spawnPositions [][]resolv.Vector
	for _, direction := range []Direction{Up, Down, Left, Right} {
		checks := CheckTilesInDirection(atPosition, direction, reach, dx, tags.TagWall, false, ecs)
		fmt.Println("dir", direction, "checks", checks)
		fmt.Println("spawnPos", spawnPositions)
		// Maybe don't trim it here, as explostion will need know
		spawnPositions = append(spawnPositions, TakeUntilNonEmpty(checks))
	}

	// Add center explosion
	spawnPositions = append(spawnPositions, []resolv.Vector{atPosition})
	// probably will need to agument the return value with the direction to
	// choose the right animation frames...? or at least rotate them
	// by 90-degrees each?
	return spawnPositions
}

func CreateExplosion(position resolv.Vector, reach int, ecs *ecs.ECS) {
	dx := GetWorldTileDiameter(ecs)
	var spawnPostions [][]resolv.Vector = GetExplosionPositions(position, reach, ecs)
	fmt.Println("Create explosion spawn pos", spawnPostions)
	for _, dir := range spawnPostions {
		for _, pos := range dir {
			fmt.Println("Explosion pos", pos)
			explosionEntry := archtypes.Explosion.Spawn(ecs)
			components.Explosion.Set(explosionEntry, &components.ExplosionData{
				Power:          reach, // this should be redundant
				CountdownTicks: GetTickCount(ecs) + 100,
			})
			// This is were we will handle the animation
			components.Sprite.Set(explosionEntry, &components.SpriteData{
				// TODO use for loop index here to choose different animation frame
				Image:     assets.Explosion.SpriteSheet,
				Animation: assets.Explosion.CenterAnimation[0],
			})
			position = SnapToGridTileCenter(pos, dx)
			bbox := resolv.NewRectangle(position.X, position.Y, dx, dx)
			bbox.Tags().Set(tags.TagExplosion)
			components.ConvexPolygonBBox.Set(explosionEntry, bbox)
		}
	}
	// finally spawn center explosion
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
		if tileShapeTags.Has(forTags) {
			tileResult.IsEmpty = false
			tileResult.CollisionObjectTags |= tileShapeTags
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
func CheckTile(checkPosition resolv.Vector, radius float64, debugMode bool, ecs *ecs.ECS) (tileShapeTags resolv.Tags, isEmpty bool) {
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
		TestAgainst: tc.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			// lets report what we touched with rather
			//playerShape.Circle.MoveVec(set.MTV)
			// set.OtherShape = The other shape involved in the contact.
			// iscontainedby is what we need here!
			insideWall := checkPosition.IsInside(set.OtherShape)
			if insideWall {
				tileShapeTags |= *set.OtherShape.Tags()
				return false // stop testing for further intersection
			}
			//set.OtherShape.(*resolv.ConvexPolygon).IsContainedBy(set.OtherShape)
			//fmt.Println("COLLISION tag", set.OtherShape)
			return true
		},
	})
	isEmpty = true
	if tileShapeTags == 0 {
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
		//screen.DrawImage(explosionSprite.Image, op)
		// Get Bomb animation
		//explosionData := components.Explosion.Get(entry)
		drawOpts := ganim8.DrawOpts(offsetX, offsetY)
		explosionSprite.Animation.Draw(screen, drawOpts)
	}
}
