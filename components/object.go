package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/solarlune/resolv"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/features/transform"
)

type RectangleData struct {
	X, Y, W, H float64
}

func NewRectangle(x, y, w, h float64) *RectangleData {
	return &RectangleData{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

// iota Enum better?
type TileID int
type tileGrid [][]TileID

type TileGridData struct {
	Grid         tileGrid
	TileDiameter float64 // of each Tile in Grid
}

func NewTileGridData(tileGrid tileGrid, tileDiameter float64) *TileGridData {
	tg := &TileGridData{
		Grid:         tileGrid,
		TileDiameter: tileDiameter,
	}
	return tg
}

type TileMapData map[TileID]*ebiten.Image

func NewTileMap() TileMapData {
	tm := TileMapData{
		//TileID(1): ebiten.NewImage(),
	}
	return tm
}

func Level1() *TileGridData {
	tg := tileGrid{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
	}

	return NewTileGridData(tg, 16.0)
}

// Fix: Level layout not intuitive..
// doesn't look like you write it in here
func LevelLa() *TileGridData {
	tg := tileGrid{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 1, 0, 1, 0},
		{0, 1, 0, 1, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 1, 0, 1, 0, 1, 0},
		{1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}

	return NewTileGridData(tg, 16.0)
}

func PrintGrid(tg tileGrid) {
	for _, row := range tg {
		for _, v := range row {
			fmt.Printf("%d ", v)
		}
		fmt.Println()
	}
	dx := 16
	for x, row := range tg {
		for y, _ := range row {
			fmt.Printf("(%f,%f) ", float64(x*dx), float64(y*dx))
		}
		fmt.Println()
	}
}

// Use SpriteData instead?
type TileData struct {
	Id TileID
}

// Put this in to a system instead on track the tiles
// type TileObject struct {
// 	Grid             tileGrid
// 	Tilesize         float64
// 	OffsetX, OffsetY float64
// 	ScaleXY          float64 // by how much to scale all tiles when rendering
// 	imageMap         map[tileID]*ebiten.Image
// }

// Represents a position on a grid, will be used
// to position things in the Arena
type GridPositionData struct {
	X, Y int
}

// indicates whether something something is walkable (tree vs bush)
type CollidableData struct {
	IsSolid bool
}

// Will represent a particular tile/sprite that we will
// draw in tilemap cells or for the player. A map of pics,
// for easy implementation
type SpriteData struct {
	Image *ebiten.Image
}

// I can give it a default value in the parenthesis here..
var TileGrid = donburi.NewComponentType[TileGridData]()
var GridPosition = donburi.NewComponentType[GridPositionData]()
var Sprite = donburi.NewComponentType[SpriteData]()
var Rectangle = donburi.NewComponentType[RectangleData]()
var Image = donburi.NewComponentType[ebiten.Image]()
var Tile = donburi.NewComponentType[TileData]()
var TileMap = donburi.NewComponentType[TileMapData]()
var Collidable = donburi.NewComponentType[CollidableData]()

// used for collision detection/response
// IShape is an interface that can fit any shape
// like ConvexRectangle and Circle which should
// be the only shapes we need

// Circle Bounding Box. Becuase can't use a generic container
// like resolv.IShape becuase the "NewComponentType" creates a
// pointer of it and at this point I can't put things inside it
var CircleBBox = donburi.NewComponentType[resolv.Circle]()
var ConvexPolygonBBox = donburi.NewComponentType[resolv.ConvexPolygon]()

// doesn't work...
//var Object = donburi.NewComponentType[resolv.IShape]()

func SetCircleBBox(circle *resolv.Circle, tf *transform.TransformData, sprite *ebiten.Image) {
	halfW := float64(sprite.Bounds().Dx() / 2)
	h := float64(sprite.Bounds().Dy())

	x := tf.LocalPosition.X
	y := tf.LocalPosition.Y

	scaleX := tf.LocalScale.X
	circle.SetRadius(halfW * scaleX)

	newX := x + halfW //*scaleX
	newY := y + h - halfW
	circle.SetPosition(newX, newY)
}

func CircleBottomLeftPos(c *resolv.Circle) resolv.Vector {
	x := c.Position().X - c.Radius()
	y := c.Position().Y - c.Radius()
	return resolv.NewVector(x, y)
}
