package render

import (
	"boids/config"
	"boids/sim"
	"fmt"
	"image/color"
	"math"

	"github.com/ebitenui/ebitenui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	World     *sim.World
	Cfg       *config.Config
	Texture   *ebiten.Image
	DebugMode bool
	ShowUI    bool
	Ui        ebitenui.UI
	uiScale   float64
}

func (g *Game) toggleDebug() {
	g.DebugMode = !g.DebugMode
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustReleased(ebiten.KeyD) {
		g.toggleDebug()
		fmt.Println(g.DebugMode)
	}
	if inpututil.IsKeyJustReleased(ebiten.KeyTab) {
		g.ShowUI = !g.ShowUI
	}
	g.World.Update()
	if g.ShowUI {
		g.Ui.Update()
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, b := range g.World.Boids {
		g.drawBoid(screen, b, g.Texture)
	}
	if g.ShowUI {
		g.Ui.Draw(screen)
	}
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %.0f", ebiten.ActualFPS()))
}

func uiScaleForWidth(width int) float64 {
	s := float64(width) / 1920.0
	s = math.Max(0.55, math.Min(1.2, s))
	return math.Round(s*20) / 20 // quantize to 0.05 steps
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if outsideWidth != g.World.Width || outsideHeight != g.World.Height {
		ratioX := float64(outsideWidth) / float64(g.World.Width)
		ratioY := float64(outsideHeight) / float64(g.World.Height)
		g.World.UpdatePositions(ratioX, ratioY)
		g.World.Width = outsideWidth
		g.World.Height = outsideHeight
	}
	if newScale := uiScaleForWidth(outsideWidth); newScale != g.uiScale {
		g.uiScale = newScale
		g.Ui = BuildUI(g.World, g.Cfg, g.uiScale)
	}
	return outsideWidth, outsideHeight
}

func (g *Game) drawBoid(screen *ebiten.Image, boid sim.Boid, texture *ebiten.Image) {
	angle := math.Atan2(boid.VY, boid.VX)
	size := boid.Size
	vertices := make([]ebiten.Vertex, 3)
	offsets := []float64{0, math.Pi * 3 / 4, -math.Pi * 3 / 4}
	for i := range 3 {
		vertices[i] = ebiten.Vertex{
			DstX:   float32(boid.X + size*math.Cos(angle+offsets[i])),
			DstY:   float32(boid.Y + size*math.Sin(angle+offsets[i])),
			ColorR: 1,
			ColorG: 1,
			ColorB: 1,
			ColorA: 1,
		}
	}

	screen.DrawTriangles(vertices, []uint16{0, 1, 2}, texture, nil)

	// debug stuff
	if g.DebugMode {
		vector.StrokeCircle(screen, float32(boid.X), float32(boid.Y), float32(g.World.DetectionRadius), 1, color.RGBA{80, 80, 80, 80}, true)
		vector.StrokeCircle(screen, float32(boid.X), float32(boid.Y), float32(g.World.AvoidanceRadius), 1, color.RGBA{0, 180, 0, 80}, true)
	}
}
