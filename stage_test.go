package main

import (
	"fmt"
	"testing"
)

func TestStage(t *testing.T) {
	fmt.Println("Building Stage")
	stg := DebugStage2()
	fmt.Println("Debug Stage Layout:")
	PrintStage(stg)
}
