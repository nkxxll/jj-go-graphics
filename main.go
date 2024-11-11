package main

import (
	"jj-go-graphics/games"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	startScreen = games.NewStartScreen()
	nextGame    = ""
	gameRunning = false
	game        games.Game
	running     = true
)

func main() {
	// create the window
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(games.TITLE, sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, games.WIDTH, games.HEIGHT, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}

	// this is the upper game loop that renders the game
	// this loop holds the control over the window
	for running {
		if gameRunning {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				game.HandleEvent(event, &running)
				if running == false {
					gameRunning = false
				}
			}
			loopTime := game.Loop(surface, &gameRunning)
			window.UpdateSurface()

			delay := (1000 / games.FRAMERATE) - loopTime
			if delay > 4_294_967_295 {
				println("Quitting..Error delay bigger than 32 bit uint")
				running = false
				continue
			}
			// NOTE: this is weird
			sdl.Delay(uint32(delay))
		} else {
			for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				nextGame = startScreen.HandleEvent(event, &running)
			}
			switch nextGame {
			case "Snake":
				game = games.NewSnakeGame()
				gameRunning = true
				nextGame = ""
			}

			loopTime := startScreen.Loop(surface, &running)
			window.UpdateSurface()

			delay := (1000 / games.FRAMERATE) - loopTime
			if delay > 4_294_967_295 {
				println("Quitting..Error delay bigger than 32 bit uint")
				running = false
				continue
			}
			// NOTE: this is weird
			sdl.Delay(uint32(delay))
		}
	}
}
