package main

import (
	"fmt"
	"testing"
)

func TestStage(t *testing.T) {
	fmt.Println("Building Stage")
	stg := DebugStage3()
	fmt.Println("Debug Stage Layout:")
	PrintStage(stg)
}
