package main

import (
	"github.com/raidancampbell/goby/cartridge"
	"github.com/raidancampbell/goby/cpu"
	"github.com/raidancampbell/goby/mem"
	"github.com/raidancampbell/goby/render"
	"os"
	"path/filepath"
)

type dmg struct {
	cpu  *cpu.CPU
	cart *cartridge.ROM
	ram  *mem.RAM
	lcd  render.LCD
	ppu  render.PPU
}

func main() {
	cwd, err := os.Getwd()
	gamedir := filepath.Join(cwd, "omitted-assets/tetris.gb")
	romFile, err := os.OpenFile(gamedir, os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}

	gb := dmg{}
	gb.cart = cartridge.Load(romFile)
	gb.cpu = cpu.Get()
	gb.ram = cpu.GetRAM()
	gb.lcd = render.LCD{}
	gb.lcd.Init()
	gb.ppu = render.PPU{}
	gb.ppu.Init()

	bootrom, err := os.OpenFile(filepath.Join(cwd, "omitted-assets/dmg_boot.bin"), os.O_RDONLY, 0)
	if err != nil {
		panic(err)
	}
	cpu.LoadBootrom(bootrom)
	cpu.GetRAM().LoadCartridge(gb.cart)
	cpu.InitPCForBootrom()
	cpu.Run()
}
