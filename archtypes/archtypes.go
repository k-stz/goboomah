package archtypes

import (
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/layers"
	"github.com/k-stz/goboomah/tags"
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
		components.Tick,     // count global Ticks here
	)

	Player = newArchetype(
		tags.Player,
		// will store speed vector for input in collision system
		components.Player,
		components.Sprite,      // maps TileIDs to ebiten.Images
		transform.Transform,    // Contains transform data
		components.CircleBBox,  // used for collision detection, OLD
		components.ShapeCircle, // Collision logic
		// why don't I need to add the AddCirecle... yet?
	)

	ArenaTile = newArchetype(
		tags.Arena,
	)

	Bomb = newArchetype(
		tags.Bomb,
		components.Sprite,
		components.Bomb,
		components.ConvexPolygonBBox,
	)

	Explosion = newArchetype(
		tags.Explosion,
		components.Sprite,
		components.Explosion,
		components.ConvexPolygonBBox,
	)

	// Used to visuallize logic in game
	DebugCircle = newArchetype(
		tags.DebugCircle,
		components.Sprite,
		components.ShapeCircle,
	)

	// Used as solid wall or breakable wall
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
