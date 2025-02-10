package systems

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
)

var rect = resolv.NewRectangle(200, 100, 32, 32)

// Here we handle all the collisions detection and response
func UpdateObjects(ecs *ecs.ECS) {
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	playerShape := components.ShapeCircle.Get(playerEntry)
	dx := GetWorldTileDiameter(ecs)

	movement := resolv.NewVector(player.Direction.X, player.Direction.Y)
	// Fill these with playerData stats, so you can pick up
	// speed enhancing items!
	maxSpd := 4.0
	friction := 0.5
	accel := 0.3 + friction

	player.Movement = player.Movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)
	playerShape.Circle.MoveVec(player.Movement)

	playerShape.Circle.IntersectionTest(resolv.IntersectionTestSettings{
		TestAgainst: playerShape.Circle.SelectTouchingCells(1).
			// Filter all solids
			FilterShapes().ByTags(tags.TagWall | tags.TagBomb),
		OnIntersect: func(set resolv.IntersectionSet) bool {
			if set.OtherShape.Tags().Has(tags.TagWall) {
				playerShape.Circle.MoveVec(set.MTV)
				//fmt.Println("COLLISION with wall, applying MTV", set.MTV)
			}
			if set.OtherShape.Tags().Has(tags.TagBomb) {
				// Only apply MTV when player is mostly outside of the bomb
				// such that when a player places a bomb he can still move on it
				magMTV := set.MTV.Magnitude()
				// ratio := magMTV / dx
				if (magMTV / dx) < 0.15 {
					playerShape.Circle.MoveVec(set.MTV.Scale(1.0))
				}
				//fmt.Println("Collision Bomb, MTV!", set.MTV, "mag", magMTV, "ratio", ratio)
			}
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

// This renders the Collision Detection Bounding Boxes and Circles
func DrawPhysics(ecs *ecs.ECS, screen *ebiten.Image) {
	spaceEntry, _ := tags.Space.First(ecs.World)
	space := components.Space.Get(spaceEntry)

	space.ForEachShape(func(shape resolv.IShape, index, maxCount int) bool {

		var drawColorCircle color.Color = color.White
		var drawDebugCircle color.RGBA = color.RGBA{255, 32, 128, 255}
		var breakableWall color.RGBA = color.RGBA{0, 50, 180, 255}

		// tags := shape.Tags()

		drawColor := color.RGBA{32, 255, 128, 255}

		switch o := shape.(type) {
		// currently Circle is just the player
		case *resolv.Circle:
			color := drawColorCircle
			if o.Tags().Has(tags.TagDebug) {
				color = drawDebugCircle
			}
			vector.StrokeCircle(screen, float32(o.Position().X), float32(o.Position().Y), float32(o.Radius()), 2, color, false)
		case *resolv.ConvexPolygon:
			if shape.Tags().Has(tags.TagBreakable) {
				drawColor = breakableWall
			}
			for _, l := range o.Lines() {
				vector.StrokeLine(screen, float32(l.Start.X), float32(l.Start.Y), float32(l.End.X), float32(l.End.Y), 2, drawColor, false)
			}
		}

		return true

	})

}
