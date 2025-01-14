package main

import (
	"embed"
	"errors"
	"fmt"
	"image"
	"log/slog"

	_ "image/png"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// implement ebiten.Game
type Game struct {
	stage Stage
}

var count int = 1000

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
		return errors.New("Countdown zero!")
	}
	return nil
}

// This is called every "Frame". Frame is a time unit
// for rendering and this depends on the display's
// refresh rate! If the monitor refresh rate is 60 [Hz]
// Draw is called 60 times per second!
//
// ebiten.Image: In Ebitengine all images, like
// images created from image files, offscreen images
// (temporary render target), and the scree are represented
// as ebiten.Image objects!
//
// screen: is the final destination of rendering!
// The window shows the final state of screen every frame
// Whenever Draw is called, screen is also reset.
// that's why the ebitenutil.DebugPrint needs to be
// called on every Draw(..) rendering

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screnHeight int) {
	return 640, 480
}

// func mustLoadImage(name string) *ebiten.Image {
// 	f, err := assets.Open(name)
// 	if err != nil {
// 		fmt.Println("cant open assets! name:", name)
// 		panic(err)
// 	}
// 	defer f.Close()

// 	//img, _, err := image.Decode(f)
// 	img, _, err := ebitenutil.NewImageFromFile(name)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//fmt.Println("eimg", eimg)
// 	//fmt.Println("img", img)
// 	return img
// }

// All files directly under the "assets/"-dir will be
// embedded in the go binary and will be accessible from
// the "assets" variable of type embed.FS

//go:embed assets/*
var assets embed.FS

func mustLoadImage(name string) *ebiten.Image {
	f, err := assets.Open(name)
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
	tilesize := g.stage.tilesize
	// op.GeoM.Scale(4, 4)
	op := &ebiten.DrawImageOptions{}
	for x, row := range g.stage.Grid {
		for y, val := range row {
			if val == 0 {
				// TODO: can't I set the translation directly
				// instead of Reset()ing it each time?
				op.GeoM.Translate(float64(y*tilesize), float64(x*tilesize))
				screen.DrawImage(meadow_tile, op)
				op.GeoM.Reset()
			}
		}
	}
	// screen.DrawImage(meadow_tile, op)
	// op.GeoM.Translate(16*4, 16*4)
	// screen.DrawImage(meadow_tile, op)
}

func (g *Game) Draw(screen *ebiten.Image) {
	count--
	move++
	screen.Fill(color.RGBA{0xff, 0, 0, 0xff})

	message := fmt.Sprintf("Hello, World\ncount: %d\n", count)
	ebitenutil.DebugPrint(screen, message)
	screen.DrawImage(gopher, nil)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(0.10, 0.10)
	op.GeoM.Translate(10, move)
	screen.DrawImage(PlayerSprite, op)

	g.DrawStage(screen)
}

var PlayerSprite = mustLoadImage("assets/player.png")
var gopher = mustLoadImage("assets/gopher.png")
var meadow_tile *ebiten.Image = mustLoadImage("assets/tiles/meadow.png")

func NewGame(tilesize int) *Game {
	debugStage := DebugStage2()
	debugStage.tilesize = tilesize
	return &Game{
		stage: debugStage,
	}
}

func main() {
	tilesize := meadow_tile.Bounds().Dx()
	game := NewGame(tilesize)
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("GoBoomer")
	if err := ebiten.RunGame(game); err != nil {
		slog.Error("ebiten.RunGame", "err", err)
	}
}
