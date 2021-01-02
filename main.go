package main

import (
	"github.com/raidancampbell/goby/cartridge"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	gamedir := filepath.Join(cwd, "omitted-assets/tetris.gb")
	file, err := os.OpenFile(gamedir, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	c := cartridge.Load(file)

	println(c.GetTitle())
}