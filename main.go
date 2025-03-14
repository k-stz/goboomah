package main

import (
	"embed"
	"fmt"
	"image"
	"log/slog"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/k-stz/goboomah/assets"
	"github.com/k-stz/goboomah/scenes"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update()
}

// implement ebiten.Game
type Game struct {
	stage Stage
	// ecs test
	bounds                    image.Rectangle
	scene                     Scene
	ScreenWidth, ScreenHeight int
}

var count int = 10000

// called every tick
// Tick is a time unit for logical updating
// the default value is 1/60 [s], then update is called
// 60 times per second by default
// i.e. the Ebitengine works in 60 ticks-per-second
// Ebitengine game game suspends when Update() returns
// non nil error!
func (g *Game) Update() error {
	count--
	// if count < 0 {
	// 	// ebiten.Termination signifies a clean, wanted, termination of the game
	// 	return ebiten.Termination
	// }
	// shutdown
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ebiten.Termination
	}

	g.scene.Update()
	return nil
}

// Draw draws the game screen.
// Draw is called every frame (typically 1/60[s] for 60Hz display).
// So roughly every 16ms, in this time everything needs to be rendered
// This is called every "Frame". Frame is a time unit
// for rendering and this depends on the display's
// refresh rate! If the monitor refresh rate is 60 [Hz]
// Draw is called 60 times per second!
//
// Screen is an ebiten.Image. In Ebitengine all images, like
// images created from image files, offscreen images
// (temporary render target), and the screen are represented
// as ebiten.Image objects!
//
// screen: is the final destination of rendering!
// The window shows the final state of screen every frame
// Whenever Draw is called, screen is also reset.
// that's why the ebitenutil.DebugPrint needs to be
// called on every Draw(..) rendering
func (g *Game) Draw(screen *ebiten.Image) {
	g.scene.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screnHeight int) {
	return 640, 480
}

// All files directly under the "assets/"-dir will be
// embedded in the go binary and will be accessible from
// the "assets" variable of type embed.FS

//go:embed assets/*
var assetsTest embed.FS

func mustLoadImage(name string) *ebiten.Image {
	fmt.Println("main.go load image name:", name)
	fmt.Println("assets embed.FS:", assetsTest)

	f, err := assetsTest.Open(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

var PlayerSprite = mustLoadImage("assets/pics/player.png")
var gopher = mustLoadImage("assets/pics/gopher.png")
var meadow_tile *ebiten.Image = mustLoadImage("assets/tiles/meadow.png")
var bush_tile *ebiten.Image = mustLoadImage("assets/tiles/bush.png")
var wall_tile *ebiten.Image = mustLoadImage("assets/tiles/wall.png")

// use it to draw the arena!
// func NewGame(tilesize float64) *Game {
// 	debugStage := DebugStage3()
// 	debugStage.tilesize = tilesize
// 	return &Game{
// 		stage: *debugStage,
// 	}
// }

func NewGame(config Config) *Game {
	assets.MustLoadAssets()
	g := &Game{
		bounds: image.Rectangle{},
		scene: &scenes.GameScene{
			ScreenWidth:  config.ScreenWidth,
			ScreenHeight: config.ScreenHeight,
			TextX:        (float64(config.ScreenWidth) / 3),
			TextY:        float64(config.ScreenHeight) / 2,
		},
		ScreenWidth:  config.ScreenWidth,
		ScreenHeight: config.ScreenHeight,
	}
	return g
}

type Config struct {
	ScreenWidth, ScreenHeight int
}

func main() {
	//tilesize := meadow_tile.Bounds().Dx()
	//gameold := NewGame(float64(tilesize))
	//game := NewEcsTestGame()
	//game := NewGame(16.0)
	config := Config{
		ScreenWidth:  480,
		ScreenHeight: 640,
	}

	game := NewGame(config)
	//game.stage = gameold.stage

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GoBoomah")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("ebiten.RunGame", "err", err)
	}
}
