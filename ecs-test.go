package main

import (
	"fmt"
	"image"
	"image/color"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
)

// ecs test
type Scene interface {
	Update()
	Draw(screen *ebiten.Image)
}

type GameScene struct {
	ecs  *ecs.ECS
	once sync.Once
}

// resolve is a 2d collision detection lib
// func GetObject(entry *donburi.Entry) *resolv.Object {
// 	return components.Object.Get(entry)
// }

var Wall = donburi.NewTag().SetName("Wall")

func DrawWall(ecs *ecs.ECS, screen *ebiten.Image) {
	//for e := range Wall.Iter(ecs.World) {

	// TODO ecs.World is empty...?
	for entry := range MyWall.Iter(ecs.World) {
		fmt.Println("Wall.Iter finally works!, entry", entry)
		// o := GetObject(e)
		// now to
		wall := MyWall.GetValue(entry)
		drawColor := color.RGBA{60, 60, 60, 255}
		ebitenutil.DrawRect(screen, wall.x, wall.y, wall.w, wall.h, drawColor)
		//ebitenutil.DrawRect(screen, o.X, o.Y, o.W, o.H, drawColor)
	}
}

const (
	Default ecs.LayerID = iota
)

func UpdateObjects(ecs *ecs.ECS) {
}

type WallComponent struct {
	x, y, w, h float64
}

var MyWall = donburi.NewComponentType[WallComponent]()

// the the ECS gets initialized
func (gs *GameScene) configure() {
	ecs := ecs.NewECS(donburi.NewWorld())

	ecs.AddSystem(UpdateObjects)

	// creates a new entity
	// component data will be initialized by default value of the struct
	myWallEntitty := ecs.World.Create(MyWall)
	entry := ecs.World.Entry(myWallEntitty)

	MyWall.SetValue(entry, WallComponent{x: 100, y: 100, w: 100.0, h: 100.0})
	ecs.AddRenderer(Default, DrawWall)
	gs.ecs = ecs
}

func (gs *GameScene) Update() {
	gs.once.Do(gs.configure)
	// runs all systems!
	// since we haven't added any yet, it should be a noop
	gs.ecs.Update()
}

func (gs *GameScene) Draw(screen *ebiten.Image) {
	//screen.Fill(color.RGBA{20, 20, 40, 255})
	gs.ecs.Draw(screen)
}

func NewGameEcs() *Game {
	g := &Game{
		bounds: image.Rectangle{},
		scene:  &GameScene{},
	}
	return g
}
