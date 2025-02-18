package systems

import (
	"fmt"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdatePlayer(ecs *ecs.ECS) {

	//tileSpiralEffect(ecs)
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerShape := components.ShapeCircle.Get(playerEntry)

	// Currently just used for damage calculation and gameover condition
	processPlayerState(playerEntry, ecs)

	var x float64
	var y float64
	dx := GetWorldTileDiameter(ecs)
	playerTilePosition := SnapToGridTileCenter(playerShape.Circle.Position(), dx)

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		x = -1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		x = 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		y = -1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		y = 1.0
	}
	player.Direction.X = x
	player.Direction.Y = y
	// so we're not quicker in any direction
	player.Direction = player.Direction.Unit()

	scaleSpeed := 0.01
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		playerShape.Scale += scaleSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		playerShape.Scale -= scaleSpeed
	}
	rotateSpeed := 0.1
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		playerShape.Rotation += rotateSpeed

	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		playerShape.Rotation -= rotateSpeed
	}

	// For explosion collision debugging
	if inpututil.IsKeyJustPressed(ebiten.KeyB) {

		pos := playerShape.Circle.Position()
		CreateExplosion(pos, 3, ecs)

		// CreateExplosion(playerTilePosition, 0, ecs)
	}

	// Bombs
	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		if CanPlaceBombs(playerTilePosition, ecs) {
			player.Bombs--
			CreateBomb(playerShape.Circle.Position(), player, ecs)
		} else {
			fmt.Println("Player bombs exhausted")
		}
	}

	// TileGrid
	if inpututil.IsKeyJustPressed(ebiten.KeyM) {
		grid := GetArenaTileGrid(ecs)
		components.PrintGrid(grid)
		// for _, row := range *grid {
		// 	for _, v := range row {
		// 		fmt.Printf("%d ", v)
		// 	}
		// 	fmt.Println()
		// }
		//g := *grid
		//e := (*grid)[0][0]
		e := grid[0][0]
		fmt.Println("Before", e)
		e = (e + 1) % 2
		grid[0][0] = e
		fmt.Println("After", e)

	}

	// For Debugging
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		//checkTiles := 4
		var tileContents []TileContent
		tileContents = CheckTilesInDirection(playerTilePosition, Down, 4, dx, tags.TagWall, true, ecs)
		for _, tc := range tileContents {
			fmt.Printf("Down Check pos %v = %v (%s)\n", tc.CenterPosition,
				tc.IsEmpty, tags.ToString(tc.CollisionObjectTags))
		}
		tileContents = CheckTilesInDirection(playerTilePosition, Up, 4, dx, tags.TagWall, true, ecs)
		for _, tc := range tileContents {
			fmt.Printf("Up Check pos %v = %v (%s)\n", tc.CenterPosition,
				tc.IsEmpty, tags.ToString(tc.CollisionObjectTags))
		}
		tileContents = CheckTilesInDirection(playerTilePosition, Right, 4, dx, tags.TagWall, true, ecs)
		for _, tc := range tileContents {
			fmt.Printf("Right Check pos %v = %v (%s)\n", tc.CenterPosition,
				tc.IsEmpty, tags.ToString(tc.CollisionObjectTags))
		}
		tileContents = CheckTilesInDirection(playerTilePosition, Left, 4, dx, tags.TagWall, true, ecs)
		for _, tc := range tileContents {
			fmt.Printf("Left Check pos %v = %v (%s)\n", tc.CenterPosition,
				tc.IsEmpty, tags.ToString(tc.CollisionObjectTags))
		}
	}
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	ticks := GetTickCount(ecs)
	for entry := range tags.Player.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		playerSprite := components.Sprite.Get(entry)

		halfW := float64(playerSprite.Image.Bounds().Dx() / 2)
		halfH := float64(playerSprite.Image.Bounds().Dy() / 2)

		circleShape := components.ShapeCircle.Get(entry)
		pos := circleShape.Circle.Bounds().Center()
		rad := circleShape.Circle.Radius()
		rotation := circleShape.Rotation
		diameter := max(halfW, halfH)
		// diameter * x = radius
		scale := rad / diameter

		var offsetX float64 = pos.X
		var offsetY float64 = pos.Y //- halfH + halfW

		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Scale(scale, scale)
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(offsetX, offsetY)
		// Draw invincibility Frames
		player := GetPlayer(ecs)
		if player.State == components.Invincible {
			color := float32(Oscillator(math.Cos, int(ticks), 1.0, 0.1, 4.0))
			op.ColorScale.Scale(color, color, color, 255)
		}
		screen.DrawImage(playerSprite.Image, op)
	}
}

func processPlayerState(entry *donburi.Entry, ecs *ecs.ECS) {
	currentTicks := GetTickCount(ecs)

	player := components.Player.Get(entry)
	circleShape := components.ShapeCircle.Get(entry)

	switch player.State {
	case components.Idle:
		fmt.Println("Idle", currentTicks)
		if player.Damaged {
			player.Lives--
			// death animation 1 seconds
			player.Duration = currentTicks + (1 * 60)
			player.State = components.Death
		}
	case components.Death:
		fmt.Println("DEATH", currentTicks)
		if player.Duration < currentTicks {
			// stop in player in his tracks
			circleShape.Circle.MoveVec(resolv.NewVector(0, 0))
			player.Direction = resolv.NewVector(0, 0)
			// invincibility ticks for 3 secondds
			player.Duration = currentTicks + (2 * 60)
			player.State = components.Invincible
			player.Damaged = false
			// Respawn!
			circleShape.Circle.SetPositionVec(player.RespawnPoint)
		}
		// Death Animation
		circleShape.Rotation += 0.2
	case components.Invincible:
		circleShape.Rotation = 0.0
		fmt.Println("Invincible", currentTicks)
		if player.Duration < currentTicks {
			// Transition to Idle State after some time
			player.State = components.Idle
			player.Damaged = false
		}
	}
}
