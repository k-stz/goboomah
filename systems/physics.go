package systems

import (
	"fmt"

	"github.com/k-stz/goboomer/components"
	"github.com/k-stz/goboomer/tags"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/features/math"
	"github.com/yohamta/donburi/features/transform"
)

var rect = resolv.NewRectangle(200, 100, 32, 32)

// Here we handle all the collisions detection and response
func UpdateObjects(ecs *ecs.ECS) {
	playerEntry, _ := tags.Player.First(ecs.World)
	player := components.Player.Get(playerEntry)
	tf := transform.Transform.Get(playerEntry)
	playerSprite := components.Sprite.Get(playerEntry)
	playerCircleBBox := components.CircleBBox.Get(playerEntry)
	components.SetCircleBBox(playerCircleBBox, tf, playerSprite.Image)

	var movePlayer math.Vec2
	//fmt.Println("player speed:", player.Speed)
	if !player.Speed.IsZero() {
		fmt.Println("player speed:", player.Speed)
		fmt.Println("lets do something! Then reset it")
		// calculate player position handle colliison here
		movePlayer = player.Speed
	}
	//center := playerCircleBBox.Position()

	//bbxy := CircleBottomLeftPos(playerCircleBBox)
	//fmt.Println("Player BBox", playerCircleBBox.Position())
	//fmt.Println("corners:", playerCircleBBox.Bounds())
	playerCircleBBox.MoveVec(resolv.NewVector(player.Speed.X, player.Speed.Y))

	tf.LocalPosition = tf.LocalPosition.Add(movePlayer)

	player.Speed.X = 0
	player.Speed.Y = 0

	// for e := range components.Object.Iter(ecs.World) {
	// 	obj := collisions.GetObject(e)

	// 	//fmt.Println("obj bounds", obj.Points)
	// 	intersection := obj.Intersection(rect)
	// 	if !intersection.IsEmpty() {
	// 		fmt.Println("They're touching! Here's the data:", intersection)
	// 	}
	// }
}
