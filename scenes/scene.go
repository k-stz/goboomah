package scenes

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/factory"
	"github.com/k-stz/goboomer/layers"
	"github.com/k-stz/goboomer/systems"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

type GameScene struct {
	ecs  *ecs.ECS
	once sync.Once
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
	ecs := ecs.NewECS(donburi.NewWorld())

	//ecs.AddSystem(systems.UpdateLevelMap)

	ecs.AddSystem(systems.UpdateArena)
	ecs.AddSystem(systems.UpdatePlayer)

	ecs.AddRenderer(layers.Default, systems.DrawArena)
	ecs.AddRenderer(layers.Default, systems.DrawPlayer)

	//ecs.AddRenderer(layers.Default, systems.DrawArenaTiles)
	// Now we create the LevelMap

	// creates a new entity
	// component data will be initialized by default value of the struct
	//myWallEntitty := ecs.World.Create(MyWall)
	//entry := ecs.World.Entry(myWallEntitty)

	gs.ecs = ecs

	arenaEntry := factory.CreateArena(gs.ecs)
	playerEntry := factory.CreatePlayer(gs.ecs)

	fmt.Println("Cerated Entries", arenaEntry, playerEntry)

	//MyWall.SetValue(entry, WallComponent{x: 100, y: 100, w: 100.0, h: 100.0})
	//ecs.AddRenderer(Default, DrawWall)
}
