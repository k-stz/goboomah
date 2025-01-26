package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdatePlayer(ecs *ecs.ECS) {
	//tileSpiralEffect(ecs)
	playerEntry, _ := tags.Player.First(ecs.World)
	tf := transform.Transform.Get(playerEntry)
	player := components.Player.Get(playerEntry)
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

	//  tf.LocalPosition + math.Vec2{x, y}
	// no this has to be updated in the colliison logic
	//tf.LocalPosition = tf.LocalPosition.Add(math.NewVec2(x, y))
	// rotation for fun

	// this can be done without colision detection
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		tf.LocalRotation += 0.25
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		tf.LocalRotation -= 0.25
	}
	scaleSpeed := 0.01
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		tf.LocalScale = tf.LocalScale.Add(math.NewVec2(
			scaleSpeed, scaleSpeed))
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		tf.LocalScale = tf.LocalScale.Add(math.NewVec2(
			-scaleSpeed, -scaleSpeed))
	}
}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Player.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		playerSprite := components.Sprite.Get(entry)

		halfW := float64(playerSprite.Image.Bounds().Dx() / 2)
		halfH := float64(playerSprite.Image.Bounds().Dy() / 2)

		// TODO: use bbox as basis to draw player not the
		// other way around
		circleBBox := components.CircleBBox.Get(entry)
		pos := circleBBox.Position()
		rad := circleBBox.Radius()
		diameter := max(halfW, halfH)
		// diameter * x = radius
		scale := rad / diameter

		var offsetX float64 = pos.X
		var offsetY float64 = pos.Y //- halfH + halfW

		//tf := transform.Transform.Get(entry)
		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		//
		op.GeoM.Translate(-halfW, -halfH)
		//op.GeoM.Rotate(tf.LocalRotation)
		op.GeoM.Scale(scale, scale)
		//op.GeoM.Translate(halfW, halfH)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(playerSprite.Image, op)

		// draw bounding box of player
		// TODO: Add toggle to control this via a Config Setting
		// for the game! Put this into a bbox-render ecs system
		playerObject := components.CircleBBox.Get(entry)
		c := playerObject.Position() // position should be center for a circle...
		radius := float32(playerObject.Radius())
		vector.DrawFilledCircle(screen, float32(c.X), float32(c.Y), radius, color.RGBA{0xff, 0, 0, 10}, false)
	}
}
