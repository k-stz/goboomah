package scenes

import (
	"fmt"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/k-stz/goboomah/collisions"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/factory"
	"github.com/k-stz/goboomah/layers"
	"github.com/k-stz/goboomah/systems"
	"github.com/k-stz/goboomah/tags"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// GameScene encapsulates the ECS, systems, and scene rendering.
type GameScene struct {
	ecs                       *ecs.ECS
	once                      sync.Once
	ScreenWidth, ScreenHeight int
}

// Update calls the ECS update. The configuration is executed only once.
func (gs *GameScene) Update() {
	gs.once.Do(gs.configure)
	gs.ecs.Update()
}

// Draw renders the scene and overlays debug information.
func (gs *GameScene) Draw(screen *ebiten.Image) {
	// Clear the screen with a dark background.
	screen.Fill(color.RGBA{20, 20, 40, 255})
	gs.ecs.Draw(screen)

	// Retrieve the player entity for debug information.
	playerEntry, ok := tags.Player.First(gs.ecs.World)
	if !ok {
		ebitenutil.DebugPrintAt(screen, "Player not found", 0, 40)
		return
	}

	playerShape := components.ShapeCircle.Get(playerEntry)
	playerData := components.Player.Get(playerEntry)

	totalBombs := 0
	message := fmt.Sprintf("TPS: %0.2f\n", ebiten.ActualTPS())
	message += fmt.Sprintf("Pos: %s\n", playerShape.Circle.Position())
	message += fmt.Sprintf("SnapTileCenter: %v\n",
		systems.SnapToGridTileCenter(playerShape.Circle.Position(), systems.GetWorldTileDiameter(gs.ecs)))
	message += fmt.Sprintf("Radius: %f\nPlayerSpeed: %f\nBombs: %d\nPower: %d\nTotalBombs: %d\nTileDiameter: %02f\n",
		playerShape.Circle.Radius(),
		playerData.Movement,
		playerData.Bombs,
		playerData.Power,
		totalBombs,
		systems.GetWorldTileDiameter(gs.ecs))
	ebitenutil.DebugPrintAt(screen, message, 0, 40)
}

// configure initializes the ECS, registers systems and renderers, and creates initial entities.
func (gs *GameScene) configure() {
	// Create a new ECS world and instance.
	world := donburi.NewWorld()
	ecsInstance := ecs.NewECS(world)

	// Register game update systems.
	// ecsInstance.AddSystem(systems.UpdateLevelMap) // Uncomment if needed.
	ecsInstance.AddSystem(systems.UpdateObjects)
	ecsInstance.AddSystem(systems.UpdateArena)
	ecsInstance.AddSystem(systems.UpdateBomb)
	ecsInstance.AddSystem(systems.UpdateExplosion)
	ecsInstance.AddSystem(systems.UpdatePlayer)
	ecsInstance.AddSystem(systems.UpdateDebugCircle)

	// Register renderers for the default layer.
	ecsInstance.AddRenderer(layers.Default, systems.DrawArena)
	ecsInstance.AddRenderer(layers.Default, systems.DrawBomb)
	ecsInstance.AddRenderer(layers.Default, systems.DrawExplosion)
	ecsInstance.AddRenderer(layers.Default, systems.DrawPlayer)
	ecsInstance.AddRenderer(layers.Default, systems.DrawPhysics)
	ecsInstance.AddRenderer(layers.Default, systems.DrawDebugCircle)
	// ecsInstance.AddRenderer(layers.Default, systems.DrawArenaTiles) // Optional.

	// Assign the configured ECS to the GameScene.
	gs.ecs = ecsInstance

	// Create the world's collision grid space based on the screen dimensions.
	spaceEntry := factory.CreateSpace(gs.ecs, gs.ScreenWidth, gs.ScreenHeight)

	// Create core game objects.
	arenaEntry := factory.CreateArena(gs.ecs)
	factory.CreateSolidTiles(gs.ecs, arenaEntry)
	playerEntry := factory.CreatePlayer(gs.ecs)
	fmt.Println("Created Entries IDs:", arenaEntry.Id(), playerEntry.Id(), spaceEntry.Id())

	// Add collision bounds for solid tiles into the space.
	addSolidTilesSpace(spaceEntry, gs.ecs)

	// Optionally, set up animations here.
	// setupAnimations()
}

// addSolidTilesSpace iterates over all tile entities and adds their bounding boxes to the collision space.
func addSolidTilesSpace(spaceEntry *donburi.Entry, ecsInstance *ecs.ECS) {
	for entry := range tags.Tile.Iter(ecsInstance.World) {
		collisions.AddConvexPolygonBBox(spaceEntry, entry)
	}
}
