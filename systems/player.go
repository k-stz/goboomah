package systems

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdatePlayer(ecs *ecs.ECS) {
	//tileSpiralEffect(ecs)

}

func DrawPlayer(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Player.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		playerSprite := components.Sprite.Get(entry)
		tf := transform.Transform.Get(entry)

		var offsetX float64
		var offsetY float64
		// yes, looking up the image in a hash will kill locality
		// causing cache misses
		offsetX = tf.LocalPosition.X
		offsetY = tf.LocalPosition.Y
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Rotate(tf.LocalRotation)
		op.GeoM.Scale(tf.LocalScale.X, tf.LocalScale.Y)
		op.GeoM.Translate(offsetX, offsetY)
		screen.DrawImage(playerSprite.Image, op)
	}
}
