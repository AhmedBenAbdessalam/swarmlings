package sim

type Ling struct {
	X, Y, VX, VY float64
	Size         float64
}

func (b *Ling) Move() {
	b.X += b.VX
	b.Y += b.VY
}

func (b *Ling) Avoid(neighbors []Ling, factor, avoidanceRadius float64) (vx float64, vy float64) {
	for _, other := range neighbors {
		dx := other.X - b.X
		dy := other.Y - b.Y
		distsq := DistanceSquared(b.X, b.Y, other.X, other.Y)
		if distsq < avoidanceRadius*avoidanceRadius {
			vx -= dx / distsq
			vy -= dy / distsq
		}
	}
	return vx * factor, vy * factor
}

func (b *Ling) Align(neighbors []Ling, factor, detectionRadius float64) (vx float64, vy float64) {
	averageVX, averageVY, count := 0.0, 0.0, 0
	for _, other := range neighbors {
		if DistanceSquared(b.X, b.Y, other.X, other.Y) < detectionRadius*detectionRadius {
			averageVX += other.VX
			averageVY += other.VY
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	averageVX /= float64(count)
	averageVY /= float64(count)
	return (averageVX - b.VX) * factor, (averageVY - b.VY) * factor
}

func (b *Ling) Gather(neighbors []Ling, factor, detectionRadius float64) (vx float64, vy float64) {
	averageX, averageY, count := 0.0, 0.0, 0
	for _, other := range neighbors {
		if DistanceSquared(b.X, b.Y, other.X, other.Y) < detectionRadius*detectionRadius {
			averageX += other.X
			averageY += other.Y
			count++
		}
	}
	if count == 0 {
		return 0, 0
	}
	averageX /= float64(count)
	averageY /= float64(count)
	return (averageX - b.X) * factor, (averageY - b.Y) * factor

}

func (b *Ling) WallAvoid(width, height, margin, strength float64) (vx, vy float64) {
	if b.X < margin {
		t := (margin - b.X) / margin
		vx += strength * t * t
	} else if d := width - b.X; d < margin {
		t := (margin - d) / margin
		vx -= strength * t * t
	}
	if b.Y < margin {
		t := (margin - b.Y) / margin
		vy += strength * t * t
	} else if d := height - b.Y; d < margin {
		t := (margin - d) / margin
		vy -= strength * t * t
	}
	return vx, vy
}

func (b *Ling) Clamp(width, height float64) {
	if b.X < 0 {
		b.X = 0
		if b.VX < 0 {
			b.VX = -b.VX
		}
	} else if b.X > width {
		b.X = width
		if b.VX > 0 {
			b.VX = -b.VX
		}
	}
	if b.Y < 0 {
		b.Y = 0
		if b.VY < 0 {
			b.VY = -b.VY
		}
	} else if b.Y > height {
		b.Y = height
		if b.VY > 0 {
			b.VY = -b.VY
		}
	}
}
