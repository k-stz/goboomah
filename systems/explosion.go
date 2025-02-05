package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/archtypes"
	"github.com/k-stz/goboomah/assets"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
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
func TakeUntilNonEmpty(tileContents []TileContent) []TileContent {
	var filteredTCs []TileContent
	fmt.Println("TakeUntilEmpty Loop")
	for _, v := range tileContents {
		fmt.Println("positions", v.CenterPosition)
		if !v.IsEmpty {
			return filteredTCs
		}
		filteredTCs = append(filteredTCs, v)
	}
	return filteredTCs
}

type ExplosionOrientation struct {
	Position     resolv.Vector
	Rotation     float64
	AnimationKey string // "center", "middle" or "end"
}

func RotationFromExplosionDirection(direction Direction) (radians float64) {
	rad := 0.0
	switch direction {
	case Up:
		rad = 1.5708 * 3
	case Right:
		rad = 0.0
	case Down:
		rad = 1.5708 * 1
	case Left:
		rad = 1.5708 * 2
	}
	return rad
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
func GetExplosionPositions(atPosition resolv.Vector, reach int, ecs *ecs.ECS) (explosionOrientations []ExplosionOrientation) {
	dx := GetWorldTileDiameter(ecs)
	var spawnPositions [][]resolv.Vector
	var explosionSpawns []ExplosionOrientation
	for _, direction := range []Direction{Up, Down, Left, Right} {
		checks := CheckTilesInDirection(atPosition, direction, reach, dx, tags.TagWall, false, ecs)
		fmt.Println("dir", direction, "checks", checks)
		fmt.Println("spawnPos", spawnPositions)
		// Maybe don't trim it here, as explostion will need know
		//spawnPositions = append(spawnPositions, TakeUntilNonEmpty(checks))
		fmt.Println("check len before: ", len(checks), "reach", reach)

		checks = TakeUntilNonEmpty(checks)
		fmt.Println("check len: ", len(checks))
		// At this point we know were in a direction how many explosion need to be
		// spawned, now we can calculate their orientation
		for i, tc := range checks {
			animationKey := "middle"
			if i == reach-1 {
				animationKey = "end"
			}
			eo := ExplosionOrientation{
				Position:     tc.CenterPosition,
				Rotation:     RotationFromExplosionDirection(direction), // 90 degree
				AnimationKey: animationKey,
			}
			explosionSpawns = append(explosionSpawns, eo)
		}
	}

	// Add center explosion
	spawnPositions = append(spawnPositions, []resolv.Vector{atPosition})
	explosionSpawns = append(explosionSpawns, ExplosionOrientation{
		Position:     atPosition,
		Rotation:     0.0,
		AnimationKey: "center",
	})
	//fmt.Println("#########", spawnPositions)
	//fmt.Println(explosionSpawns)

	// probably will need to agument the return value with the direction to
	// choose the right animation frames...? or at least rotate them
	// by 90-degrees each?
	return explosionSpawns
}

func CreateExplosion(position resolv.Vector, reach int, ecs *ecs.ECS) {
	dx := GetWorldTileDiameter(ecs)
	var spawns []ExplosionOrientation = GetExplosionPositions(position, reach, ecs)
	fmt.Println("Create explosion spawn pos", position)
	for _, spawn := range spawns {
		pos := spawn.Position
		fmt.Println("Explosion pos", pos)
		explosionEntry := archtypes.Explosion.Spawn(ecs)
		components.Explosion.Set(explosionEntry, &components.ExplosionData{
			Power:          reach,                   // this should be redundant
			CountdownTicks: GetTickCount(ecs) + 500, // 50 is ideal
		})
		// This is were we will handle the animation
		components.Sprite.Set(explosionEntry, &components.SpriteData{
			// TODO use for loop index here to choose different animation frame
			Image: assets.ExplosionAnimation.SpriteSheet,
			// this shares the explosion, it is better to pass a copy here
			//Animation: assets.ExplosionAnimation.Map[spawn.AnimationKey].Clone(),
			Animation: assets.ExplosionAnimation.Map["center"].Clone(),
		})
		position = SnapToGridTileCenter(pos, dx)
		bbox := resolv.NewRectangle(position.X, position.Y, dx, dx)
		bbox.Rotate(spawn.Rotation)
		bbox.Tags().Set(tags.TagExplosion)
		components.ConvexPolygonBBox.Set(explosionEntry, bbox)
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

var ExplosionAngle float64 = 0.0

func UpdateExplosion(ecs *ecs.ECS) {
	ExplosionAngle += 1 * math.Pi / 180
	currentGameTick := GetTickCount(ecs)
	for entry := range tags.Explosion.Iter(ecs.World) {

		// SPrite update
		explosionSprite := components.Sprite.Get(entry)
		explosionSprite.Animation.Update()

		explosion := components.Explosion.Get(entry)
		if explosion.CountdownTicks <= currentGameTick {
			//fmt.Println("Blowing up!", entry.Entity())
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
		dx := GetWorldTileDiameter(ecs)
		//aniPos := SnapToGridTileTopLeft(pos, dx)
		aniPos := SnapToGridTileCenter(pos, dx)

		tileWidth, tileHeight := 48.0, 44.0
		//drawOpts := ganim8.DrawOpts(aniPos.X, aniPos.Y, 0.0, dx/tileWidth, dx/tileHeight, 0.0, 0.0)
		fmt.Println("Rotation: ", rotation)
		// THis works for 180.0 flip... very strange
		// pi := math.Pi
		//drawOpts := ganim8.DrawOpts(aniPos.X, aniPos.Y, pi/2, dx/tileWidth, dx/tileHeight, 0.0, 1.0)

		// this works
		drawOpts := ganim8.DrawOpts(aniPos.X, aniPos.Y, ExplosionAngle, dx/tileWidth, dx/tileHeight, 0.5, 0.5)

		//explosionSprite.Animation.Sprite().FlipV()
		//drawOpts.Rotate = 1
		explosionSprite.Animation.Draw(screen, drawOpts)
	}
}
