package scenes

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/collisions"
	"github.com/k-stz/goboomer/factory"
	"github.com/k-stz/goboomer/layers"
	"github.com/k-stz/goboomer/systems"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type GameScene struct {
	ecs                      *ecs.ECS
	once                     sync.Once
	ScreenWidth, ScreenHeigh int
}

func (gs *GameScene) Update() {
	gs.once.Do(gs.configure)
	gs.ecs.Update()
}

func (gs *GameScene) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{20, 20, 40, 255})
	gs.ecs.Draw(screen)
}

// the the ECS gets initialized
func (gs *GameScene) configure() {
	// add width, heigh to gamescene
	ecs := ecs.NewECS(donburi.NewWorld())

	//ecs.AddSystem(systems.UpdateLevelMap)

	ecs.AddSystem(systems.UpdateObjects)
	ecs.AddSystem(systems.UpdateArena)
	ecs.AddSystem(systems.UpdatePlayer)

	ecs.AddRenderer(layers.Default, systems.DrawArena)
	ecs.AddRenderer(layers.Default, systems.DrawPlayer)
	ecs.AddRenderer(layers.Default, systems.DrawPhysics)

	//ecs.AddRenderer(layers.Default, systems.DrawArenaTiles)
	// Now we create the LevelMap

	// creates a new entity
	// component data will be initialized by default value of the struct
	//myWallEntitty := ecs.World.Create(MyWall)
	//entry := ecs.World.Entry(myWallEntitty)

	gs.ecs = ecs

	//MyWall.SetValue(entry, WallComponent{x: 100, y: 100, w: 100.0, h: 100.0})
	//ecs.AddRenderer(Default, DrawWall)

	// Define the world's Space. Here, a Space is essentially a grid
	// (the game's width and height, or 640x360), made up of 16x16
	// cells. Each cell can have 0 or more Objects within it,
	// and collisions can be found by checking the Space to see if
	// the Cells at specific positions contain (or would contain)
	// Objects. This is a broad, simplified approach to collision
	// detection.
	spaceEntry := factory.CreateSpace(gs.ecs, gs.ScreenWidth, gs.ScreenHeigh)

	// Create objects
	arenaEntry := factory.CreateArena(gs.ecs)
	factory.CreateSolidTiles(gs.ecs, arenaEntry)
	playerEntry := factory.CreatePlayer(gs.ecs)
	fmt.Println("Created Entries IDs:", arenaEntry.Id(), playerEntry.Id(), spaceEntry.Id())

	collisions.AddCircleBBox(spaceEntry, playerEntry)
	addSolidTilesSpace(spaceEntry, ecs)
}

func addSolidTilesSpace(spaceEntry *donburi.Entry, ecs *ecs.ECS) {
	// just save them in the arenaentry instead?
	// they should be recalculated with the areana?
	// We might want to implement a Camera later for larger stages...
	for entry := range tags.Tile.Iter(ecs.World) {
		collisions.AddConvexPolygonBBox(spaceEntry, entry)
	}
}
