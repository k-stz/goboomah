package archtypes

import (
	"fmt"

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
		components.Rectangle, // offset used as origin for arena
		components.TileMap,   // maps TileIDs to ebiten.Images
		transform.Transform,  // Testing: used instead of Rectangle... 
	)

	ArenaTile = newArchetype(
		tags.Arena,
	)

	Tile = newArchetype(
		tags.Tile,
		components.GridPosition,
		components.Sprite,
		components.Collidable,
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
	fmt.Println("archtype components:", a.components)
	fmt.Println("additional components:", cs)
	e := ecs.World.Entry(ecs.Create(
		layers.Default,
		append(a.components, cs...)...,
	))
	fmt.Println("POST SPAWN")
	return e
}
