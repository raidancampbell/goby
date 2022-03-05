package render

import (
	"github.com/veandco/go-sdl2/sdl"
	"os"
	"os/signal"
)

type LCD struct {
	window *sdl.Window
}

func (l *LCD) destruct() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<- c
	l.window.Destroy()
	sdl.Quit()
}

func (l *LCD) Init() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}

	// render the whole 256x256 internal framebuffer
	window, err := sdl.CreateWindow("goby", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 256, 256, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	l.window = window

	go l.destruct()

	surface, err := window.GetSurface()
	if err != nil {
		panic(err)
	}
	surface.FillRect(nil, 0)

	window.UpdateSurface()
}