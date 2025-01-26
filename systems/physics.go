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
	"github.com/yohamta/donburi/features/math"
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
	tf := transform.Transform.Get(playerEntry)
	//playerSprite := components.Sprite.Get(playerEntry)
	playerCircle := components.ShapeCircle.Get(playerEntry).Circle
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

	movement := resolv.NewVectorZero()
	maxSpd := 4.0
	friction := 0.5
	accel := 0.5 + friction

	Movement := resolv.NewVector(player.Speed.X, player.Speed.Y)

	Movement = Movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)

	playerCircle.MoveVec(resolv.NewVector(player.Speed.X, player.Speed.Y))

	playerCircle.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: playerCircle.SelectTouchingCells(1).FilterShapes(),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			playerCircle.MoveVec(set.MTV)
			fmt.Println("COLLISION, applying MTV", set.MTV)
			// also update the tf.LocalTransform
			//player.Speed.X = 0
			//player.Speed.Y = 0
			return true
		},
	})
	movePlayer := math.NewVec2(Movement.X, Movement.Y)

	//center := playerCircleBBox.Position()

	//bbxy := CircleBottomLeftPos(playerCircleBBox)
	//fmt.Println("Player BBox", playerCircleBBox.Position())
	//fmt.Println("corners:", playerCircleBBox.Bounds())
	//playerCircleBBox.MoveVec(resolv.NewVector(player.Speed.X, player.Speed.Y))

	tf.LocalPosition = tf.LocalPosition.Add(movePlayer)

	// for e := range components.Object.Iter(ecs.World) {
	// 	obj := collisions.GetObject(e)

	// 	//fmt.Println("obj bounds", obj.Points)
	// 	intersection := obj.Intersection(rect)
	// 	if !intersection.IsEmpty() {
	// 		fmt.Println("They're touching! Here's the data:", intersection)
	// 	}
	// }
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
