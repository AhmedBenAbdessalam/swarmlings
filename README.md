# Swarmlings

Flocking sim in Go with Ebitengine. The end goal is little critters with neural net brains that evolve through natural selection, rendered with GPU shaders.

## Roadmap

- [x] Basic flocking (separation, alignment, cohesion) — 500 lings at 60 FPS
- [x] Triangle rendering with debug overlays for detection/avoidance radii
- [x] In-game UI for tweaking parameters, saved to `config.json`
- [ ] Spatial partitioning (grid or quadtree) to get past the O(n^2) neighbor check
- [ ] Assets (sprites, sounds, etc.)
- [ ] Ecosystem: food, energy, death, reproduction
- [ ] Predators with predator-prey dynamics
- [ ] Evolvable genomes with mutation and natural selection
- [ ] Save/load, schema versioning, structured logging, stress testing
- [ ] WASM build deployed to GitHub Pages
- [ ] Kage shaders for GPU-accelerated rendering and visual effects
- [ ] Neural net brains evolved through genetic algorithms
- [ ] Chaos testing and performance ceiling analysis

## Tech Stack

Go, Ebitengine v2, EbitenUI

## Project Structure

```
main.go              entry point
config.json          simulation parameters
sim/
  ling.go            ling struct + flocking rules
  sim.go             world update loop
  sim_test.go        tests
render/
  render.go          ebitengine game loop, drawing
  ui.go              parameter tuning panel
config/
  config.go          json config load/save
```

## Running

```bash
go run .
```

**Tab** toggles the parameter UI, **D** toggles debug mode (shows radii).

## Configuration

Edit `config.json` directly or use the in-game sliders:

```json
{
  "avoidance_factor": 0.8,
  "alignment_factor": 0.006,
  "gathering_factor": 0.001,
  "avoidance_radius": 10,
  "detection_radius": 120
}
```
