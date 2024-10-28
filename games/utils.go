package games

import "github.com/veandco/go-sdl2/sdl"

const (
	TITLE     = "MyGames"
	WIDTH     = 800
	HEIGHT    = 600
	FRAMERATE = 60
)

type Point struct {
	X int32
	Y int32
}

type Game interface {
	Loop(*sdl.Surface, *bool) uint64
	HandleEvent(sdl.Event, *bool) string
}

func NewPoint(x, y int32) Point {
	return Point{X: x, Y: y}
}
