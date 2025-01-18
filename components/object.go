package components

import (
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

// I can give it a default value in the parenthesis here..
var Object = donburi.NewComponentType[Rectangle]()
