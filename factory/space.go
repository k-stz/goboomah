package factory

import (
	"github.com/k-stz/goboomah/archtypes"
	"github.com/k-stz/goboomah/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateSpace(ecs *ecs.ECS, screenWidth, screenHeight int) *donburi.Entry {
	space := archtypes.Space.Spawn(ecs)

	// Space size had to be increased, else solid tiles at the edge
	// would be fully ignored from collision detection logic!
	// They could still be rendered but they wouldn't show
	// up when quering them over neighbouring cells!
	spaceData := resolv.NewSpace(screenWidth*2, screenHeight*2, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
