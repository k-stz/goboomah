package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/yohamta/donburi/ecs"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdatePlayer(ecs *ecs.ECS) {

	//tileSpiralEffect(ecs)
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerShape := components.ShapeCircle.Get(playerEntry)
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

	if ebiten.IsKeyPressed(ebiten.KeySpace) {

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

	// For Debugging
	if ebiten.IsKeyPressed(ebiten.KeyT) {
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
	for entry := range tags.Player.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		playerSprite := components.Sprite.Get(entry)

		halfW := float64(playerSprite.Image.Bounds().Dx() / 2)
		halfH := float64(playerSprite.Image.Bounds().Dy() / 2)

		circleShape := components.ShapeCircle.Get(entry)
		pos := circleShape.Circle.Position()
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
		screen.DrawImage(playerSprite.Image, op)
	}
}
