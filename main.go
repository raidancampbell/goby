package main

import (
	"github.com/raidancampbell/goby/cartridge"
	"github.com/raidancampbell/goby/cpu"
	"os"
	"path/filepath"
)

func main() {
	cwd, err := os.Getwd()
	gamedir := filepath.Join(cwd, "omitted-assets/tetris.gb")
	romFile, err := os.OpenFile(gamedir, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	cart := cartridge.Load(romFile)
	bootrom, err := os.OpenFile(filepath.Join(cwd, "omitted-assets/dmg_boot.bin"), os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	cpu.LoadBootrom(bootrom)
	cpu.GetRAM().LoadCartridge(cart)
	cpu.InitPCForBootrom()
	cpu.Run()
}