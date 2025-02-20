package assets

import (
	"embed"
	"fmt"
	"image"
	_ "image/png"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/yohamta/ganim8/v2"
)

// here we will embed the assets and export it as a package
// then it will be loaded into the ecs in scene.go!

var (
	//go:embed *
	assetsFS embed.FS

	Meadow_tile *ebiten.Image
	Blob_npc    *ebiten.Image
	Wall_tile   *ebiten.Image
	Bush_tile   *ebiten.Image

	Player             *ebiten.Image
	Bomb_tile          *ebiten.Image
	ExplosionAnimation *Animation

	FiraSansRegularSource *text.GoTextFaceSource
)

type Animation struct {
	SpriteSheet *ebiten.Image
	// so you can use Animation[walkRight], Animation[centerExplosion] etc
	Map map[string]*ganim8.Animation
}

// Load spritesheet in assets/animation folder. You just have
// to refer to it by name
func NewExplosionAnimation(spritFileName string) *Animation {
	spritesheet := mustLoadImage("animation/" + spritFileName)
	//imageWidth := spritesheet.Bounds().Dx()
	//imageHeight := spritesheet.Bounds().Dy()
	imageWidth := 336
	imageHeight := 134

	// explosion seems to be 336x134 pixels
	// with 3 rows of 7 frames each
	// so  336/7 = 48 width per frame
	// and 134/3 = 44.6 => 45 height per frame...
	// strange, they do look like same size
	frameWidth := 48  //48
	frameHeight := 48 // 45
	// you can specify different time duration for each frame like this:
	// []time.Duration { 100 * time.Millisecond, 100 * time.Millisecond }
	m150 := 150 * time.Millisecond
	m100 := 100 * time.Millisecond
	//Center
	gCenter := ganim8.NewGrid(frameWidth, frameHeight, imageWidth, imageHeight, 0, 0, 0)
	fmt.Println("gridCenter:", gCenter)
	centerExplosionAnimation := ganim8.New(spritesheet, gCenter.GetFrames("1-7", 1),
		[]time.Duration{m100, m100, m150, m150, m150, m100, m100})
	// Sides
	frameWidth = 48
	frameHeight = 43
	gSides := ganim8.NewGrid(frameWidth, frameHeight, imageWidth, imageHeight, 0, 48, 0)
	fmt.Println("gridSide:", gSides)
	middleExplosionAnimation := ganim8.New(spritesheet, gSides.GetFrames("1-7", 1),
		[]time.Duration{m100, m100, m150, m150, m150, m100, m100})
	endExplosionAnimation := ganim8.New(spritesheet, gSides.GetFrames("1-7", 2),
		[]time.Duration{m100, m100, m150, m150, m150, m100, m100})

	animation := &Animation{
		SpriteSheet: spritesheet,
		Map: map[string]*ganim8.Animation{
			"center": centerExplosionAnimation,
			"middle": middleExplosionAnimation,
			"end":    endExplosionAnimation,
		},
	}
	//	animation.Map["center"] = centerExplosionAnimation
	return animation
}

func MustLoadAssets() {
	Meadow_tile = mustLoadImage("tiles/meadow.png")
	Wall_tile = mustLoadImage("tiles/wall.png")
	Bush_tile = mustLoadImage("tiles/bush.png")
	Player = mustLoadImage("pics/gopher.png")
	Bomb_tile = mustLoadImage("tiles/bomb.png")
	Blob_npc = mustLoadImage("pics/blob.png")
	ExplosionAnimation = NewExplosionAnimation("explosion.png")
	FiraSansRegularSource = mustLoadFont("fonts/FiraSans-Regular.ttf")
}

func mustLoadFont(name string) *text.GoTextFaceSource {
	fmt.Println("assetsfs:", assetsFS)
	f, err := assetsFS.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	source, err := text.NewGoTextFaceSource(f)
	if err != nil {
		panic(err)
	}

	return source
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
