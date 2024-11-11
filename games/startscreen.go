package games

import (
	"github.com/veandco/go-sdl2/sdl"
)

var (
	white            = sdl.Color{R: 255, G: 255, B: 255, A: 255}
	textMargineRight = 20
	colors           = []sdl.Color{
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 255, G: 255, B: 0, A: 255},   // Yellow
		{R: 0, G: 255, B: 255, A: 255},   // Cyan
		{R: 255, G: 0, B: 255, A: 255},   // Magenta
		{R: 128, G: 128, B: 128, A: 255}, // Gray
		{R: 255, G: 165, B: 0, A: 255},   // Orange
		{R: 128, G: 0, B: 128, A: 255},   // Purple
		{R: 139, G: 69, B: 19, A: 255},   // Brown
	}
)

type Menu struct {
	idx  int
	name string
}

type StartScreen struct {
	menues     []Menu
	tiles      int
	tileHeight int
}

func (s StartScreen) Loop(surface *sdl.Surface, running *bool) (loopTime uint64) {

	startTime := sdl.GetTicks64()
	s.DrawMenu(surface)

	// Calculate time passed since start of the function
	endTime := sdl.GetTicks64()
	return endTime - startTime
}

func (s StartScreen) DrawMenu(surface *sdl.Surface) {
	for idx, _ := range s.menues {
		x := 0
		y := idx * s.tileHeight
		rect := sdl.Rect{X: int32(x), Y: int32(y), W: WIDTH, H: int32(s.tileHeight)}
		color := colors[idx%len(colors)]
		pixel := sdl.MapRGBA(surface.Format, color.R, color.G, color.B, color.A)
		surface.FillRect(&rect, pixel)
	}
}

func (s StartScreen) HandleEvent(event sdl.Event, running *bool) string {
	switch t := event.(type) {
	case *sdl.QuitEvent:
		*running = false
		return ""
	case *sdl.MouseButtonEvent:
		if t.State == sdl.RELEASED {
			y := t.Y
			for idx := range s.tiles {
				if y < int32((idx+1)*s.tileHeight) {
					print(s.menues[idx].name)
					return s.menues[idx].name
				}
			}
		}
		return ""
	}
	return ""
}

func NewStartScreen() StartScreen {
	menues := []Menu{NewMenu(0, "Snake"), NewMenu(1, "Planet"), NewMenu(2, "Space")}
	tiles := len(menues)
	return StartScreen{
		menues:     menues,
		tiles:      tiles,
		tileHeight: HEIGHT / tiles,
	}
}

func NewMenu(idx int, name string) Menu {
	return Menu{
		idx:  idx,
		name: name,
	}
}
