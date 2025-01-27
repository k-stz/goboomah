package assets

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

// here we will embed the assets and export it as a package
// then it will be loaded into the ecs in scene.go!

var (
	//go:embed *
	assetsFS embed.FS

	Meadow_tile *ebiten.Image
	Wall_tile   *ebiten.Image
	Player      *ebiten.Image
	Bomb_tile   *ebiten.Image
)

func MustLoadAssets() {
	Meadow_tile = mustLoadImage("tiles/meadow.png")
	Wall_tile = mustLoadImage("tiles/wall.png")
	Player = mustLoadImage("large/gopher.png")
	Bomb_tile = mustLoadImage("tiles/bomb.png")
}

func mustLoadImage(name string) *ebiten.Image {
	fmt.Println("assetsfs:", assetsFS)
	f, err := assetsFS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		// "must" function naming convention so it panics
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}
