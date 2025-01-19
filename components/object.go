package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

type Rectangle struct {
	X, Y, W, H float64
}

func NewRectangle(x, y, w, h float64) *Rectangle {
	return &Rectangle{
		X: x,
		Y: y,
		W: w,
		H: h,
	}
}

type Player struct {
	health int
	badges int
}

// iota Enum better?
type TileID int
type tileGrid [][]int

type TileGridData struct {
	Grid tileGrid
}

func NewTileGridData(tileGrid tileGrid) *TileGridData {
	tg := &TileGridData{
		Grid: tileGrid,
	}
	return tg
}

func Level1() *TileGridData {
	tg := tileGrid{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
	}

	return NewTileGridData(tg)
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
	// TODO implement a map[string]image
	SpriteID string
}

// I can give it a default value in the parenthesis here..
var Object = donburi.NewComponentType[Rectangle]()
var Image = donburi.NewComponentType[ebiten.Image]()
var Tile = donburi.NewComponentType[TileData]()
var TileGrid = donburi.NewComponentType[TileGridData]()
var GridPosition = donburi.NewComponentType[GridPositionData]()
var Sprite = donburi.NewComponentType[SpriteData]()
var Collidable = donburi.NewComponentType[SpriteData]()
