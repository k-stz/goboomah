package main

import "fmt"

type cell struct {
	solid bool
}

type tilemap [][]int

type Stage struct {
	Grid             tilemap
	tilesize         float64
	offsetX, offsetY float64
	scaleXY          float64 // by how much to scale all tiles when rendering
}

func NewStage(grid tilemap, tilesize, offsetX, offsetY, scaleXY float64) *Stage {
	return &Stage{
		Grid:     grid,
		tilesize: tilesize,
		offsetX:  offsetX,
		offsetY:  offsetY,
		scaleXY:  scaleXY,
	}

}

func PrintStage(stage *Stage) {
	for _, row := range stage.Grid {
		for _, v := range row {
			fmt.Printf("%d ", v)
		}
		fmt.Println()
	}
	dx := 16
	for x, row := range stage.Grid {
		for y, _ := range row {
			fmt.Printf("(%f,%f) ", float64(x*dx), float64(y*dx))
		}
		fmt.Println()
	}
}

func DebugStage3() *Stage {
	a := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
	}
	stage := NewStage(a, 1.0, 100.0, 100.0, 4.0)
	return stage
}
