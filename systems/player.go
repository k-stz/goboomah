package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
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
	var x float64
	var y float64
	velocity := 1.0
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		x = -velocity
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		x = velocity
	}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		y = -velocity
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		y = velocity
	}
	//  tf.LocalPosition + math.Vec2{x, y}
	tf.LocalPosition = tf.LocalPosition.Add(math.NewVec2(x, y))
	// rotation for fun
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

		tf := transform.Transform.Get(entry)
		var offsetX float64
		var offsetY float64
		// yes, looking up the image in a hash will kill locality
		// causing cache misses
		offsetX = tf.LocalPosition.X
		offsetY = tf.LocalPosition.Y
		op := &ebiten.DrawImageOptions{}
		// translate to origin, so scaling and rotation work
		// intuitively
		op.GeoM.Translate(-halfW, -halfH)
		op.GeoM.Rotate(tf.LocalRotation)
		op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
		op.GeoM.Translate(halfW, halfH)

		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(playerSprite.Image, op)
	}
}
