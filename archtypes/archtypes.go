package archtypes

import (
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/layers"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

var (
	Arena = newArchetype(
		tags.Arena,
		components.TileGrid,
		components.TileMap,  // maps TileIDs to ebiten.Images
		transform.Transform, // Testing: used instead of Rectangle...
	)

	Player = newArchetype(
		tags.Player,
		// will store speed vector for input in collision system
		components.Player,
		components.Sprite,     // maps TileIDs to ebiten.Images
		transform.Transform,   // Contains transform data
		components.CircleBBox, // used for collision detection
		// why don't I need to add the AddCirecle... yet?
	)

	ArenaTile = newArchetype(
		tags.Arena,
	)

	Tile = newArchetype(
		tags.Tile,
		//components.GridPosition,
		components.Sprite,
		components.ConvexPolygonBBox,
		components.Collidable,
	)

	// "resolv"'s physics engine "space"
	Space = newArchetype(
		tags.Space,
		components.Space,
	)
	// FloatingPlatform = newArchetype(
	// 	tags.FloatingPlatform,
	// 	components.Object,
	// 	components.Tween,
	// )
	// Player = newArchetype(
	// 	tags.Player,
	// 	components.Player,
	// 	components.Object,
	// )
	// Ramp = newArchetype(
	// 	tags.Ramp,
	// 	components.Object,
	// )
	// Space = newArchetype(
	// 	components.Space,
	// )
	// Wall = newArchetype(
	// 	tags.Wall,
	// 	components.Object,
	// )
)

// type Arena struct {

// 	// Tiles      [][]TileComponent // Grid of tiles
// 	// Obstacles  []Entity          // Static obstacles (e.g., walls, blocks)
// 	// Walkable   []Entity          // Walkable spaces
// 	// Bounds     Rectangle         // Arena boundaries
// }

type archetype struct {
	components []donburi.IComponentType
}

func newArchetype(cs ...donburi.IComponentType) *archetype {
	return &archetype{
		components: cs,
	}
}

// Spawn adds a new entry to the ecs
func (a *archetype) Spawn(ecs *ecs.ECS, cs ...donburi.IComponentType) *donburi.Entry {
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.components, cs...)...,
	))
	return e
}
