package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/transform"
)

var rect = resolv.NewRectangle(200, 100, 32, 32)

func UpdateObjects2(ecs *ecs.ECS) {
	playerEntry, _ := tags.Player.First(ecs.World)
	//player := components.Player.Get(playerEntry)
	tf := transform.Transform.Get(playerEntry)
	playerSprite := components.Sprite.Get(playerEntry)
	playerCircleBBox := components.CircleBBox.Get(playerEntry)
	components.SetCircleBBox(playerCircleBBox, tf, playerSprite.Image)

	c := playerCircleBBox

	movement := resolv.NewVectorZero()
	maxSpd := 4.0
	friction := 0.5
	accel := 0.5 + friction

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		movement.X -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		movement.X += 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		movement.Y -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		movement.Y += 1
	}

	Movement := resolv.NewVector(tf.LocalPosition.X, tf.LocalPosition.Y)
	Movement = Movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)

	//tf.LocalPosition = math.NewVec2(Movement.X, Movement.Y)

	c.MoveVec(Movement)

	c.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: c.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			fmt.Println("Intersectiong!")
			c.MoveVec(set.MTV)
			return true
		},
	})
}

// Here we handle all the collisions detection and response
func UpdateObjects(ecs *ecs.ECS) {
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	//tf := transform.Transform.Get(playerEntry)
	//playerSprite := components.Sprite.Get(playerEntry)
	playerShape := components.ShapeCircle.Get(playerEntry)
	//components.SetCircleBBox(playerCircleBBox, tf, playerSprite.Image)

	//fmt.Println("player speed:", player.Speed)
	// if player.Speed.IsZero() {
	// 	return
	// }
	// calculate player position handle colliison here
	// Collision Detection and Response time!

	count := 0
	for entry := range components.Tile.Iter(ecs.World) {
		count++
		fmt.Println("Tiles exists, yes?", count, entry.Id())
	}

	// This doesn't modulate the speed properly
	movement := resolv.NewVector(player.Speed.X, player.Speed.Y)
	maxSpd := 1.0
	friction := 0.5
	accel := 0.5 + friction

	movement = movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)

	playerShape.Circle.MoveVec(movement.Scale(5))

	playerShape.Circle.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: playerShape.Circle.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			playerShape.Circle.MoveVec(set.MTV)
			fmt.Println("COLLISION, applying MTV", set.MTV)
			// also update the tf.LocalTransform
			//player.Speed.X = 0
			//player.Speed.Y = 0
			return true
		},
	})

	// Update scale
	circle := playerShape.Circle
	circle.SetRadius(playerShape.Scale * playerShape.Radius)
}

func DrawPhysics(ecs *ecs.ECS, screen *ebiten.Image) {
	spaceEntry, _ := tags.Space.First(ecs.World)
	space := components.Space.Get(spaceEntry)

	space.ForEachShape(func(shape resolv.IShape, index, maxCount int) bool {

		var drawColorCircle color.Color = color.White

		// tags := shape.Tags()

		drawColor := color.RGBA{32, 255, 128, 255}

		switch o := shape.(type) {
		case *resolv.Circle:
			vector.StrokeCircle(screen, float32(o.Position().X), float32(o.Position().Y), float32(o.Radius()), 2, drawColorCircle, false)
		case *resolv.ConvexPolygon:

			for _, l := range o.Lines() {
				vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, drawColor, false)
			}
		}

		return true

	})

}
