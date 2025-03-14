package systems

import (
	"fmt"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/components"
	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"golang.org/x/exp/rand"
)

// Remove Bomb from ecs and its object from the Collisoin spacce
func RemoveEnemy(enemyEntry *donburi.Entry, ecs *ecs.ECS) {
	space := GetSpace(ecs)
	boundingCircle := components.ShapeCircle.Get(enemyEntry).Circle
	space.Remove(boundingCircle)
	ecs.World.Remove(enemyEntry.Entity())
}

// This function is called when the player was caught in an
// explosion
func processEnemyExplosion(enemyEntry *donburi.Entry, ecs *ecs.ECS) {
	currentTicks := GetTickCount(ecs)
	state := components.Explodable.Get(enemyEntry)
	circleShape := components.ShapeCircle.Get(enemyEntry)
	aiState := components.AI.Get(enemyEntry)
	if !state.ProcessedExplosion && state.Exploding {
		//enemyStat.Lives--
		state.ProcessedExplosion = true
		aiState.Hp--
	}
	// Enemy exploding state: set
	if !state.Despawn && state.Exploding && state.ExplodingTick <= currentTicks {
		state.Exploding = false
		state.DespawnTick = currentTicks + (60 * 2)
		fmt.Println("Enemy Hp", aiState.Hp)
		if aiState.Hp <= 0 {
			fmt.Println("I'm despawning", currentTicks, "despawn=", state.Despawn)
			state.Despawn = true
		}
	}
	// Enemy hurt state
	if state.Exploding && state.ExplodingTick > currentTicks {
		// Play damage animation
		circleShape.Rotation -= 0.1
		fmt.Println("Enemy: Im exploding!")
		fmt.Println("I'm exploding", currentTicks, "despawn=", state.Despawn)
	} else {
		circleShape.Rotation = 0
	}
	//Enemy dying state
	if state.Despawn {
		circleShape.Rotation += 0.1
	}
	//Enemy dead state, despawn for good
	if state.Despawn && state.DespawnTick <= currentTicks {
		fmt.Println("removing enemy", currentTicks)
		// blob enemy blows up on death!
		// probably track for which enemy to do this
		CreateExplosion(circleShape.Circle.Position(), 2, ecs)
		RemoveEnemy(enemyEntry, ecs)
	}

}

// Returns a vector pointing a a cardinal direction at random
// (1,0), (0,1), (-1,0), (0,-1)
func randomCardinalDirection() resolv.Vector {
	directions := []resolv.Vector{
		{X: 0, Y: -1}, // North
		{X: 1, Y: 0},  // East
		{X: 0, Y: 1},  // South
		{X: -1, Y: 0}, // West
	}

	rand.Seed(uint64(time.Now().UnixNano()))
	return directions[rand.Intn(len(directions))]
}

func processEnemyState(enemyEntry *donburi.Entry, ecs *ecs.ECS) {
	currentTicks := GetTickCount(ecs)
	//explodeState := components.Explodable.Get(enemyEntry)
	circleShape := components.ShapeCircle.Get(enemyEntry)
	aiState := components.AI.Get(enemyEntry)

	switch aiState.State {
	case components.Idle:
		chosenDirection := randomCardinalDirection()
		aiState.Direction = chosenDirection
		// transition to walking state for some time
		aiState.State = components.Walking
		aiState.Duration = currentTicks + 20
	case components.Walking:
		if aiState.Duration < currentTicks {
			// Transition to Idle State for some time
			aiState.State = components.Idle
			aiState.Duration = currentTicks + 100
		}
	}
	// Apply movement to enemy

	movement := resolv.NewVector(aiState.Direction.X, aiState.Direction.Y)

	baseSpeed := 0.1
	aiState.Movement = aiState.Movement.Scale(baseSpeed)
	maxSpd := 1.0
	friction := 0.5
	accel := 0.01 * friction
	aiState.Movement = aiState.Movement.Add(movement.Scale(accel)).SubMagnitude(friction).ClampMagnitude(maxSpd)
	aiState.Movement = aiState.Movement.Add(movement)
	circleShape.Circle.MoveVec(aiState.Movement)

}

// Here handle player input and update velocity/movement
// laeter collision response will be collected from here
func UpdateEnemy(ecs *ecs.ECS) {
	for enemyEntry := range tags.Enemy.Iter(ecs.World) {
		processEnemyState(enemyEntry, ecs)
		processEnemyExplosion(enemyEntry, ecs)
	}
}

// CosineOscillator calculates a cosine wave value that oscillates between two values
// ticks: current game tick count
// Used for creating a cycled squish animation for the blob enemy
// using one oscilation to control the x-axis scaling and
// and another one for the y-axis in the opposite direction
// period: duration in seconds for a full cycle
// from: the minimum oscillation value
// to: the maximum oscillation value
func Oscillator(fn func(float64) float64, ticks int, period float64, from float64, to float64) float64 {
	const ticksPerSecond = 60.0
	frequency := 1.0 / period // Hz (cycles per second)
	angle := 2 * math.Pi * frequency * (float64(ticks) / ticksPerSecond)
	scaledCosine := (math.Cos(angle) + 1) / 2 // Scale between 0 and 1
	return from + (to-from)*scaledCosine      // Scale to the desired range
}

func DrawEnemy(ecs *ecs.ECS, screen *ebiten.Image) {
	for entry := range tags.Enemy.Iter(ecs.World) {
		//o := dresolv.GetObject(e)
		enemySprite := components.Sprite.Get(entry)
		state := components.Explodable.Get(entry)

		halfW := float64(enemySprite.Image.Bounds().Dx() / 2)
		halfH := float64(enemySprite.Image.Bounds().Dy() / 2)

		circleShape := components.ShapeCircle.Get(entry)

		// pos should be the center of the circle
		pos := circleShape.Circle.Position()
		//rad := circleShape.Circle.Radius
		//scale := circleShape.Scale
		rotation := circleShape.Rotation
		//diameter := max(halfW, halfH)
		// diameter * x = radius
		//scale := rad / diameter

		var offsetX float64 = pos.X
		var offsetY float64 = pos.Y //- halfH + halfW

		op := &ebiten.DrawImageOptions{}

		// translate to origin, so scaling and rotation work
		// intuitively
		ticks := GetTickCount(ecs)
		// make them squish not in sync to be more natural?
		scaleX := Oscillator(math.Cos, int(ticks), 3.0, 1.5, 0.6)
		scaleY := Oscillator(math.Sin, int(ticks), 3.0, 0.6, 1.5)

		if state.Exploding {
			op.ColorScale.Scale(float32(scaleX), 2.0, 2.0, float32(scaleY))
		}

		op.GeoM.Translate(-halfW, -halfH)
		if state.Despawn {
			color := float32(Oscillator(math.Cos, int(ticks), 1.0, 0.1, 4.0))
			op.ColorScale.Scale(color, color, color, 255)
			circleShape.Scale += 0.01
			scaleX = circleShape.Scale
			scaleY = circleShape.Scale
		}
		op.GeoM.Scale(scaleX, scaleY)
		//op.GeoM.Scale(1.0, 1.0) // undo for debugging
		op.GeoM.Rotate(rotation)
		op.GeoM.Translate(offsetX, offsetY)
		//fmt.Println("draw enemmy at", offsetX, offsetY)
		screen.DrawImage(enemySprite.Image, op)
	}
}
