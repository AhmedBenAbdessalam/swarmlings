package main

import (
	"swarmlings/config"
	"swarmlings/render"
	"swarmlings/sim"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// initialize dots
	boids := make([]sim.Boid, 500)
	world := sim.New(boids, 800, 600)
	for i := range world.Boids {
		world.Boids[i] = sim.Boid{
			X:    rand.Float64() * float64(world.Width),
			Y:    rand.Float64() * float64(world.Height),
			VX:   rand.Float64() * 1,
			VY:   rand.Float64() * 1,
			Size: 10,
		}
	}

	cfg := config.Load()

	// apply config to world
	world.AvoidanceFactor = cfg.AvoidanceFactor
	world.AlignmentFactor = cfg.AlignmentFactor
	world.GatheringFactor = cfg.GatheringFactor
	world.AvoidanceRadius = cfg.AvoidanceRadius
	world.DetectionRadius = cfg.DetectionRadius
	world.MaxSpeed = cfg.MaxSpeed

	ui := render.BuildUI(&world, &cfg, 1.0)

	// create game instance
	texture := ebiten.NewImage(1, 1)
	texture.Fill(color.White)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	game := &render.Game{World: &world, Cfg: &cfg, Texture: texture, ShowUI: true, Ui: ui}
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
