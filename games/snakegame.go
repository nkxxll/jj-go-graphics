package games

import (
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type SnakeBody struct {
	bodySlice []Point
	Len       int
}

type SnakeGame struct {
	apples    []Point
	appleLen  int
	Direction Point
	Body      SnakeBody
}

func (b SnakeBody) GetHead() Point {
	return b.bodySlice[0]
}

func (b SnakeBody) DrawBody(surface *sdl.Surface) {
	for idx := range b.Len {
		// Draw on the surface
		p := b.bodySlice[idx]
		rect := sdl.Rect{X: p.X - 5, Y: p.Y - 5, W: 10, H: 10}
		colour := sdl.Color{R: 255, G: 0, B: 255, A: 255} // purple
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
}

func (b *SnakeBody) Move(direction Point) {
	old := b.bodySlice[0]
	b.bodySlice[0].X += direction.X * 5
	b.bodySlice[0].Y += direction.Y * 5
	for i := 1; i < b.Len; i++ {
		tmp := b.bodySlice[i]
		b.bodySlice[i] = old
		old = tmp
	}
}

func (g *SnakeGame) CheckApple() {
	// check for eaten apple
	head := g.Body.GetHead()
	for idx := 0; idx < g.appleLen; idx++ {
		apple := g.apples[idx]
		if g.RectCollision(apple, head, 10) {
			g.EatApple(idx)
		}
	}
}

func (b *SnakeBody) AddOne(direction Point) {
	oldTail := b.bodySlice[b.Len-1]
	b.Move(direction)
	b.Move(direction)
	b.Len += 1
	b.bodySlice[b.Len-1] = oldTail
}

func NewSnakeBody(x, y int32) SnakeBody {
	slice := make([]Point, 1024)
	slice[0] = Point{x, y}
	return SnakeBody{
		bodySlice: slice,
		Len:       1,
	}
}

func NewSnakeGame() SnakeGame {
	applesSlice := make([]Point, 256)

	// spawn first apple at the start
	randx := rand.Intn(WIDTH-50) + 50
	randy := rand.Intn(HEIGHT-50) + 50
	applesSlice[0] = Point{X: int32(randx), Y: int32(randy)}

	playerPos := NewPoint((WIDTH / 2), int32(HEIGHT/2))

	return SnakeGame{
		apples:    applesSlice,
		appleLen:  1,
		Direction: NewPoint(0, 0),
		Body:      NewSnakeBody(playerPos.X, playerPos.Y),
	}
}

func (g *SnakeGame) SpawnApple() {
	// too much apples
	if g.appleLen == 256 {
		return
	}
	randx := rand.Intn(WIDTH-50) + 50
	randy := rand.Intn(HEIGHT-50) + 50

	g.appleLen += 1
	g.apples[g.appleLen-1] = Point{int32(randx), int32(randy)}
}

func (g SnakeGame) DrawApples(surface *sdl.Surface) {
	for idx := range g.appleLen {
		apple := g.apples[idx]
		rect := sdl.Rect{X: apple.X - 5, Y: apple.Y - 5, W: 10, H: 10}
		colour := sdl.Color{R: 0, G: 255, B: 0, A: 255}
		pixel := sdl.MapRGBA(surface.Format, colour.R, colour.G, colour.B, colour.A)
		surface.FillRect(&rect, pixel)
	}
}

func (g *SnakeGame) EatApple(idx int) {
	g.appleLen -= 1
	if g.appleLen > 0 {
		last := g.apples[g.appleLen]
		g.apples[idx] = last
	}
	g.Body.AddOne(g.Direction)
}

func (g *SnakeGame) HandleEvent(event sdl.Event, running *bool) string {
	switch t := event.(type) {
	case *sdl.QuitEvent: // NOTE: Please use `*sdl.QuitEvent` for `v0.4.x` (current version).
		println("Quitting..")
		*running = false
		break
	case *sdl.KeyboardEvent:
		if t.State == sdl.RELEASED {
			if t.Keysym.Sym == sdl.K_LEFT || t.Keysym.Sym == sdl.K_h {
				if g.Direction.X == 1 {
					return ""
				}
				g.Direction.X = -1
				g.Direction.Y = 0
			} else if t.Keysym.Sym == sdl.K_RIGHT || t.Keysym.Sym == sdl.K_l {
				if g.Direction.X == -1 {
					return ""
				}
				g.Direction.X = 1
				g.Direction.Y = 0
			}
			if t.Keysym.Sym == sdl.K_UP || t.Keysym.Sym == sdl.K_k {
				if g.Direction.Y == 1 {
					return ""
				}
				g.Direction.X = 0
				g.Direction.Y = -1
			} else if t.Keysym.Sym == sdl.K_DOWN || t.Keysym.Sym == sdl.K_j {
				if g.Direction.Y == -1 {
					return ""
				}
				g.Direction.X = 0
				g.Direction.Y = 1
			}
		}
		break
	}
	return ""
}

func (g *SnakeGame) CheckBodyCollision() bool {
	// check for collision with own body
	head := g.Body.GetHead()
	for idx := 1; idx < g.Body.Len; idx++ {
		body := g.Body.bodySlice[idx]
		if g.RectCollision(head, body, 5) {
			println("Game Over...")
			return true
		}
	}
	return false
}

// Loop runs the game loop sets running to false if the game is lost
func (g *SnakeGame) Loop(surface *sdl.Surface, running *bool) (loopTime uint64) {
	// Get time at the start of the function
	startTime := sdl.GetTicks64()
	if rand.Float32() < 0.01 {
		g.SpawnApple()
	}

	// Update player position
	g.Body.Move(g.Direction)

	head := g.Body.GetHead()
	if head.X < 0 {
		println("game over..")
		*running = false
		return 0
	} else if head.X > WIDTH {
		println("game over..")
		*running = false
		return 0
	}
	if head.Y < 0 {
		println("game over..")
		*running = false
		return 0
	} else if head.Y > HEIGHT {
		println("game over..")
		*running = false
		return 0
	}

	collision := g.CheckBodyCollision()
	if collision {
		*running = false
		return 0
	}

	g.CheckApple()

	// Clear surface
	surface.FillRect(nil, 0)

	g.Body.DrawBody(surface)
	g.DrawApples(surface)

	// Calculate time passed since start of the function
	endTime := sdl.GetTicks64()
	return endTime - startTime
}

func (g SnakeGame) RectCollision(a, b Point, r int) bool {
	half_r := int32(r / 2)
	upperLeftA := Point{X: a.X - half_r, Y: a.Y - half_r}
	lowerRightA := Point{X: a.X + half_r, Y: a.Y + half_r}
	upperLeftB := Point{X: b.X - half_r, Y: b.Y - half_r}
	lowerRightB := Point{X: b.X + half_r, Y: b.Y + half_r}
	if upperLeftA.X > lowerRightB.X {
		return false
	}
	if upperLeftA.Y > lowerRightB.Y {
		return false
	}
	if lowerRightA.X >= upperLeftB.X && lowerRightA.Y >= upperLeftB.Y {
		return true
	}
	return false
}
