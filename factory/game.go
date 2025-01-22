package factory

import (
	"github.com/k-stz/goboomer/archtypes"
	"github.com/k-stz/goboomer/components"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

func CreateGame(ecs *ecs.ECS, screenWidth, screenHeight int) *donburi.Entry {
	space := archtypes.Space.Spawn(ecs)

	spaceData := resolv.NewSpace(screenWidth, screenHeight, 16, 16)
	components.Space.Set(space, spaceData)

	return space
}
