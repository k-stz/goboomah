package components

import (
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
)

type State int

const (
	Idle State = iota
	Walking
	Chasing
	Dodging
)

type AIState struct {
	// represents in which direction and with what magnitude
	// the AI wants to move
	// the collision/physics system will then calculate what
	// position this results into and will update the
	// Monsters actual position (stored in another component)
	Direction resolv.Vector
	// Current movement Speed of the AI Monster, this is by how much
	// the Monster will move in a direction
	Movement resolv.Vector
	Hp       int
	State    State
	// used to track for how long to stay in the next state
	Duration TickCount
}

// Create AI with given hp (healthpoits).
func NewAI(hp int) *AIState {
	return &AIState{
		Hp:    hp,
		State: Idle,
	}
}

var AI = donburi.NewComponentType[AIState]()
