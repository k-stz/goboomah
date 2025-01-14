package main

import "fmt"

type cell struct {
	solid bool
}

type Stage struct {
	Layout [][]int
}

func NewStage(layout [][]int) {

}

func PrintStage(stage Stage) {
	for _, row := range stage.Layout {
		for _, v := range row {
			fmt.Printf("%d ", v)
		}
		fmt.Println()
	}
}

func DebugStage2() Stage {
	a := [][]int{
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 0, 0},
	}
	return Stage{Layout: a}
}

func DebugStage() Stage {
	row0 := []int{0, 0, 0, 0, 0}
	row1 := []int{0, 1, 0, 1, 0}
	row2 := []int{0, 0, 0, 0, 0}
	row3 := []int{0, 1, 0, 1, 0}
	row4 := []int{0, 0, 0, 0, 0}

	stage := Stage{
		Layout: [][]int{
			row0,
			row1,
			row2,
			row3,
			row4,
		},
	}
	fmt.Println("stage:", stage)
	return stage
}
