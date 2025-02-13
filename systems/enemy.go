package systems

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/yohamta/donburi/ecs"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdateEnemy(ecs *ecs.ECS) {

	// //tileSpiralEffect(ecs)
	// enemyEntry, _ := tags.Enemy.First(ecs.World)
	// enemyShape := components.ShapeCircle.Get(enemyEntry)

	// scaleSpeed := 0.01
	// if ebiten.IsKeyPressed(ebiten.KeyG) {
	// 	enemyShape.Scale += scaleSpeed

	// }
	// if ebiten.IsKeyPressed(ebiten.KeyF) {
	// 	enemyShape.Scale -= scaleSpeed
	// }
	// rotateSpeed := 0.1
	// if ebiten.IsKeyPressed(ebiten.KeyR) {
	// 	enemyShape.Rotation += rotateSpeed

	// }
	// if ebiten.IsKeyPressed(ebiten.KeyE) {
	// 	enemyShape.Rotation -= rotateSpeed
	// }
}

func DrawEnemy(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Enemy.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		enemySprite := components.Sprite.Get(entry)

		halfW := float64(enemySprite.Image.Bounds().Dx() / 2)
		halfH := float64(enemySprite.Image.Bounds().Dy() / 2)

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
		op.GeoM.Scale(1.0, 1.0) // undo for debugging
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(offsetX, offsetY)
		fmt.Println("draw enemmy at", offsetX, offsetY)
		screen.DrawImage(enemySprite.Image, op)
	}
}
