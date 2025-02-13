package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/donburi"
)

// iota Enum better?
type TileID int
type Grid [][]TileID

type TileGridData struct {
	Grid         Grid
	TileDiameter float64 // of each Tile in Grid
}

func NewTileGridData(tileGrid Grid, tileDiameter float64) *TileGridData {
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
	tg := Grid{
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
	tg := Grid{
		{0, 0, 0, 2, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 1, 0},
		{0, 2, 0, 2, 1, 0, 0, 0, 2, 2},
		{0, 0, 0, 0, 1, 0, 2, 0, 2, 0},
		{0, 0, 0, 0, 0, 0, 0, 2, 1, 0},
		{0, 0, 1, 1, 1, 0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0, 0, 1, 0, 1, 0},
		{0, 1, 9, 1, 0, 2, 0, 2, 2, 2},
		{0, 9, 1, 1, 1, 1, 2, 2, 1, 2},
		{1, 0, 0, 0, 0, 0, 2, 2, 2, 2},
	}

	return NewTileGridData(tg, 16.0)
}

func LevelTile1() *TileGridData {
	tg := Grid{
		{1},
	}
	return NewTileGridData(tg, 16.0)
}

func PrintGrid(tg Grid) {
	for _, row := range tg {
		for _, v := range row {
			fmt.Printf("%d ", v)
		}
		fmt.Println()
	}
	// dx := 16
	// for x, row := range tg {
	// 	for y, _ := range row {
	// 		fmt.Printf("(%f,%f) ", float64(x*dx), float64(y*dx))
	// 	}
	// 	fmt.Println()
	// }
}

// Use SpriteData instead?
type TileData struct {
	// the offset in the tilemap the tile belongs to
	GridX, GridY int
}

// Represents a position on a grid, will be used
// to position things in the Arena
type GridPositionData struct {
	X, Y int
}

// I can give it a default value in the parenthesis here..
var TileGrid = donburi.NewComponentType[TileGridData]()
var GridPosition = donburi.NewComponentType[GridPositionData]()
var Tile = donburi.NewComponentType[TileData]()
var TileMap = donburi.NewComponentType[TileMapData]()
