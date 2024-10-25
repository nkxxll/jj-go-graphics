package main

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	TITLE     = "MyGames"
	WIDTH     = 800
	HEIGHT    = 600
	FRAMERATE = 60
)

type Point struct {
	x int32
	y int32
}

type Body struct {
	bodySlice []Point
	len       int
}

type Game struct {
	apples   []Point
	appleLen int
}

func (b Body) DrawBody(surface *sdl.Surface) {
	for idx := range b.len {
		// Draw on the surface
		p := b.bodySlice[idx]
		rect := sdl.Rect{X: p.x - 5, Y: p.y - 5, W: 10, H: 10}
		colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
}

func (b *Body) Move(direction Point) {
	old := b.bodySlice[0]
	b.bodySlice[0].x += direction.x * 5
	b.bodySlice[0].y += direction.y * 5
	for i := 1; i < b.len; i++ {
		tmp := b.bodySlice[i]
		b.bodySlice[i] = old
		old = tmp
	}
}

func (b *Body) AddOne(direction Point) {
	oldTail := b.bodySlice[b.len-1]
	b.Move(direction)
	b.Move(direction)
	b.len += 1
	b.bodySlice[b.len-1] = oldTail
}

func NewBody(x, y int32) Body {
	slice := make([]Point, 1024)
	slice[0] = Point{x, y}
	return Body{
		bodySlice: slice,
		len:       1,
	}
}

func NewGame() Game {
	applesSlice := make([]Point, 256)

	// spawn first apple at the start
	randx := rand.Intn(WIDTH-50) + 50
	randy := rand.Intn(HEIGHT-50) + 50
	applesSlice[0] = Point{x: int32(randx), y: int32(randy)}

	return Game{
		apples:   applesSlice,
		appleLen: 1,
	}
}

func (g *Game) SpawnApple() {
	// too much apples
	if g.appleLen == 256 {
		return
	}
	randx := rand.Intn(WIDTH-50) + 50
	randy := rand.Intn(HEIGHT-50) + 50

	g.appleLen += 1
	g.apples[g.appleLen-1] = Point{int32(randx), int32(randy)}
}

func (g Game) DrawApples(surface *sdl.Surface) {
	for idx := range g.appleLen {
		apple := g.apples[idx]
		rect := sdl.Rect{X: apple.x - 5, Y: apple.y - 5, W: 10, H: 10}
		colour := sdl.Color{R: 0, G: 255, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
}

func (g *Game) EatApple(idx int, body *Body, direction Point) {
	g.appleLen -= 1
	if g.appleLen > 0 {
		last := g.apples[g.appleLen]
		g.apples[idx] = last
	}
	body.AddOne(direction)
}

func (g Game) RectCollision(a, b Point, r int) bool {
	half_r := int32(r / 2)
	upperLeftA := Point{x: a.x - half_r, y: a.y - half_r}
	lowerRightA := Point{x: a.x + half_r, y: a.y + half_r}
	upperLeftB := Point{x: b.x - half_r, y: b.y - half_r}
	lowerRightB := Point{x: b.x + half_r, y: b.y + half_r}
	if upperLeftA.x > lowerRightB.x {
		return false
	}
	if upperLeftA.y > lowerRightB.y {
		return false
	}
	if lowerRightA.x >= upperLeftB.x && lowerRightA.y >= upperLeftB.y {
		return true
	}
	return false
}

var (
	playerX, playerY = int32(WIDTH / 2), int32(HEIGHT / 2)
	// playerX, playerY = int32(0), int32(0)
	body      = NewBody(playerX, playerY)
	game      = NewGame()
	direction = Point{x: int32(0), y: int32(0)}
	running   = true
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			handleEvent(event)
		}

		loopTime := loop(surface)
		window.UpdateSurface()

		delay := (1000 / FRAMERATE) - loopTime
		if delay > 4_294_967_295 {
			println("Quitting..Error delay bigger than 32 bit uint")
			running = false
			continue
		}
		// NOTE: this is weird
		sdl.Delay(uint32(delay))
	}
}

func handleEvent(event sdl.Event) {
	switch t := event.(type) {
	case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
		println("Quitting..")
		running = false
		break
	case *sdl.KeyboardEvent:
		if t.State == sdl.RELEASED {
			if t.Keysym.Sym == sdl.K_LEFT || t.Keysym.Sym == sdl.K_h {
				direction.x = -1
				direction.y = 0
			} else if t.Keysym.Sym == sdl.K_RIGHT || t.Keysym.Sym == sdl.K_l {
				direction.x = 1
				direction.y = 0
			}
			if t.Keysym.Sym == sdl.K_UP || t.Keysym.Sym == sdl.K_k {
				direction.x = 0
				direction.y = -1
			} else if t.Keysym.Sym == sdl.K_DOWN || t.Keysym.Sym == sdl.K_j {
				direction.x = 0
				direction.y = 1
			}
		}
		break
	}
}

func loop(surface *sdl.Surface) (loopTime uint64) {
	// Get time at the start of the function
	startTime := sdl.GetTicks64()
	if rand.Float32() < 0.01 {
		game.SpawnApple()
	}

	// Update player position
	body.Move(direction)

	head := body.bodySlice[0]
	if head.x < 0 {
		println("game over..")
		running = false
		return 0
	} else if head.x > WIDTH {
		println("game over..")
		running = false
		return 0
	}
	if head.y < 0 {
		println("game over..")
		running = false
		return 0
	} else if head.y > HEIGHT {
		println("game over..")
		running = false
		return 0
	}

	// check for collision with own body
	for idx := 1; idx < body.len; idx++ {
		if game.RectCollision(body.bodySlice[idx], head, 5) {
			println("Game Over...")
			running = false
			return 0
		}
	}

	// check for eaten apple
	for idx := 0; idx < game.appleLen; idx++ {
		apple := game.apples[idx]
		if game.RectCollision(apple, head, 10) {
			game.EatApple(idx, &body, direction)
		}
	}

	// Clear surface
	surface.FillRect(nil, 0)

	body.DrawBody(surface)
	game.DrawApples(surface)

	// Calculate time passed since start of the function
	endTime := sdl.GetTicks64()
	return endTime - startTime
}
