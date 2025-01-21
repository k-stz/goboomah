package main

import (
	"embed"
	"fmt"
	"image"
	"log/slog"

	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/k-stz/goboomer/assets"
	"github.com/k-stz/goboomer/scenes"
)

type Scene interface {
	Draw(screen *ebiten.Image)
	Update()
}

// implement ebiten.Game
type Game struct {
	stage Stage
	// ecs test
	bounds                   image.Rectangle
	scene                    Scene
	ScreenWidth, ScreenHeigh int
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
	if count < 0 {
		// ebiten.Termination signifies a clean, wanted, termination of the game
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
	count--
	move++
	//screen.Fill(color.RGBA{0xff, 0, 0, 0xff})

	message := fmt.Sprintf("Hello, World\ncount: %d\n", count)
	ebitenutil.DebugPrint(screen, message)
	//screen.DrawImage(gopher, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.10, 0.10)
	op.GeoM.Translate(350, move)
	screen.DrawImage(PlayerSprite, op)
	g.DrawStage(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screnHeight int) {
	fmt.Println("Layout called", outsideWidth, outsideHeight)
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

var move float64 = 0

func (g *Game) DrawStage(screen *ebiten.Image) {
	scale := g.stage.scaleXY
	tilesize := g.stage.tilesize * scale
	offsetOp := &ebiten.DrawImageOptions{}
	offsetOp.GeoM.Translate(g.stage.offsetX, g.stage.offsetY)
	stageOffset := offsetOp.GeoM
	// op.GeoM.Scale(4, 4)
	op := &ebiten.DrawImageOptions{}
	for x, row := range g.stage.Grid {
		for y, val := range row {
			op.GeoM.Scale(scale, scale)
			op.GeoM.Translate(float64(y)*tilesize, float64(x)*tilesize)
			// op.GeoM.Translate(100.0, 100.0)
			op.GeoM.Concat(stageOffset)

			switch val {
			case 0:
				screen.DrawImage(meadow_tile, op)
			case 1:
				screen.DrawImage(wall_tile, op)
			}
			// TODO: can't I set the translation directly
			// instead of Reset()ing it each time?
			op.GeoM.Reset()
		}
	}
	// screen.DrawImage(meadow_tile, op)
	// op.GeoM.Translate(16*4, 16*4)
	// screen.DrawImage(meadow_tile, op)
}

var PlayerSprite = mustLoadImage("assets/large/player.png")
var gopher = mustLoadImage("assets/large/gopher.png")
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
		bounds:      image.Rectangle{},
		scene:       &scenes.GameScene{},
		ScreenWidth: config.ScreenWidth,
		ScreenHeigh: config.ScreenHeight,
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
	ebiten.SetWindowTitle("GoBoomer")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game); err != nil {
		slog.Error("ebiten.RunGame", "err", err)
	}
}
