package main

import (
	"fmt"
	"testing"

	"github.com/k-stz/goboomah/tags"
	"github.com/solarlune/resolv"
)

func TestTagsToString(t *testing.T) {
	fmt.Println("Testing tags")
	fmt.Println("tostring :", tags.ToString(tags.TagBomb|tags.TagDebug|resolv.NewTag("NewTagHere")))
}
