package systems

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

var rect = resolv.NewRectangle(200, 100, 32, 32)

// set explodable components to true for durationTicks
// where 60 ticks is roughly a second
func SetExploding(entry *donburi.Entry, ecs *ecs.ECS, durationTicks components.TickCount) {
	currentTicks := GetTickCount(ecs)
	explosionState := components.Explodable.Get(entry)
	// Check if we are already exploding, so as not
	// to reset the explosion DelayTicks
	if !explosionState.Exploding {
		explosionState.ProcessedExplosion = false
		explosionState.ExplodingTick = currentTicks + durationTicks
		explosionState.Exploding = true
	}
}

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

	collisionTags := tags.TagWall | tags.TagBomb | tags.TagEnemy | tags.TagPlayer | tags.TagExplosion
	// Entry can be any archtype that has the tag.CircleCollision
	// which currently is the Player and Enemy
	for entry := range tags.CircleCollision.Iter(ecs.World) {
		circleShape := components.ShapeCircle.Get(entry)
		circleShape.Circle.IntersectionTest(resolv.IntersectionTestSettings{
			TestAgainst: circleShape.Circle.SelectTouchingCells(1).
				// Filter all solids
				FilterShapes().ByTags(collisionTags),
			OnIntersect: func(set resolv.IntersectionSet) bool {
				if set.OtherShape.Tags().Has(tags.TagWall) {
					circleShape.Circle.MoveVec(set.MTV)
				}
				if set.OtherShape.Tags().Has(tags.TagBomb) {
					// Only apply MTV when player is mostly outside of the bomb
					// such that when a player places a bomb he can still move on it
					magMTV := set.MTV.Magnitude()
					// ratio := magMTV / dx
					if (magMTV / dx) < 0.15 {
						circleShape.Circle.MoveVec(set.MTV.Scale(1.0))
					}
				}
				// Player special logic here
				if circleShape.Circle.Tags().Has(tags.TagPlayer) {
					if set.OtherShape.Tags().Has(tags.TagEnemy) {
						fmt.Println("Collided with enemy! Ouch")
					}
					if set.OtherShape.Tags().Has(tags.TagExplosion) {
						SetExploding(entry, ecs, 60)
					}
					// explosions etcs
				}
				if circleShape.Circle.Tags().Has(tags.TagEnemy) {
					if set.OtherShape.Tags().Has(tags.TagExplosion) {
						SetExploding(entry, ecs, 60)
					}
					// special enemy logic here
					// gives football/push effect: player
					// can push the enemy around
					//circleShape.Circle.MoveVec(player.Movement)
				}
				return true
			},
		})
	}

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
		var drawDebugCircleEnemy color.RGBA = color.RGBA{32, 32, 255, 255}

		var breakableWall color.RGBA = color.RGBA{0, 50, 180, 255}

		// tags := shape.Tags()

		drawColor := color.RGBA{32, 255, 128, 255}

		switch o := shape.(type) {
		// Players and enemies have use a bounding Circle
		case *resolv.Circle:
			color := drawColorCircle
			if o.Tags().Has(tags.TagDebug) {
				color = drawDebugCircle
			}
			if o.Tags().Has(tags.TagEnemy | tags.TagDebug) {
				color = drawDebugCircleEnemy
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
