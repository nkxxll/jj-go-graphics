package games

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	MIDDLE_X = WIDTH / 2
	MIDDLE_Y = HEIGHT / 2
)

type PlanetSym struct {
	planets []Planet
}

type circle struct {
	rad   int
	color sdl.Color
}

type Planet struct {
	shape circle
	rad   int
	vel   float64
	x     int
	y     int
}

func (p Planet) Draw(surface *sdl.Surface) {
	rect := sdl.Rect{
		X: int32(p.x) - int32(p.shape.rad),
		Y: int32(p.y) - int32(p.shape.rad),
		W: int32(p.shape.rad) * 2,
		H: int32(p.shape.rad) * 2,
	}
	surface.FillRect(&rect, p.shape.color.Uint32())
}

func (p *Planet) Update() {
	alpha := math.Atan(float64(p.y) / float64(p.x))

	alpha += p.vel
	// fixme: bug here this does not work
	p.x += int(float64(p.rad) * math.Cos(alpha))
	p.y += int(float64(p.rad) * math.Sin(alpha))
}

func newCircle(rad int, color sdl.Color) circle {
	return circle{
		rad:   rad,
		color: color,
	}
}

func NewPlanet(name string, rad int, vel float64, shape circle) Planet {
	return Planet{
		shape: shape,
		rad:   rad,
		vel:   vel,
		x:     MIDDLE_X,
		y:     MIDDLE_Y - rad,
	}
}

func NewPlanetSym() *PlanetSym {
	// planets
	sunShape := newCircle(20, sdl.Color{R: 0, G: 255, B: 255, A: 200})
	earthShape := newCircle(10, sdl.Color{R: 0, G: 0, B: 255, A: 200})
	sun := NewPlanet("Sun", 0, 0, sunShape)
	earth := NewPlanet("Earth", 50, 20, earthShape)
	return &PlanetSym{
		planets: []Planet{sun, earth},
	}
}

func (g *PlanetSym) HandleEvent(event sdl.Event, running *bool) {
	switch event.(type) {
	case *sdl.QuitEvent:
		*running = false
		break
	}
}

func (g *PlanetSym) Loop(surface *sdl.Surface, running *bool) (loopTime uint64) {
	startTime := sdl.GetTicks64()

	surface.FillRect(nil, 0)

	for _, p := range g.planets {
		p.Update()
		p.Draw(surface)
	}

	// Calculate time passed since start of the function
	endTime := sdl.GetTicks64()
	return endTime - startTime

}
