package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdatePlayer(ecs *ecs.ECS) {
	// NEXT TODO: implement within which Arena Cell the player
	// is. This will be used for placing, spawns and placing
	// bombs

	//tileSpiralEffect(ecs)
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerShape := components.ShapeCircle.Get(playerEntry)
	var x float64
	var y float64
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = 1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -1.0
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = 1.0
	}
	player.Speed.X = x
	player.Speed.Y = y

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

	// Bombs
	if ebiten.IsKeyPressed(ebiten.KeySpace) {
		if CanPlaceBombs(player) {
			player.Bombs--
			CreateBomb(playerShape.Circle.Position(), player, ecs)
		} else {
			fmt.Println("Player bombs exhausted")
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		dx := GetWorldTileDiameter(ecs)
		pos := SnapToGridTileCenter(playerShape.Circle.Position(), dx)
		checkTiles := 1
		fmt.Printf("Checking %d tiles above player: ", checkTiles)
		checks := []string{}
		for i := range checkTiles {
			offsetY := -(dx + (float64(i) * dx))
			checkPos := pos.Add(resolv.NewVector(0, offsetY))
			CreateDebugCircle(checkPos, dx/2, ecs)
			fmt.Println("checking pos: ", checkPos)
			shapeTags := CheckTile(checkPos, ecs)
			tags := ""
			for _, st := range shapeTags {
				tags += st.String() + ","
			}
			checks = append(checks, tags)

		}
		fmt.Println(checks)
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
