package cartridge

import (
	"github.com/raidancampbell/goby/mem"
	"io/ioutil"
	"os"
)
const(
	minROMSize = 0x8000
	maxROMSize = 0x200000
)
type ROM []byte

func (r *ROM) LoadToRAM(m *mem.RAM) {
	for i := 0; i < 0x8000; i++ {
		m[0x100 + i] = (*r)[i]
	}
}

func Load(f *os.File) *ROM {
	var r ROM
	b, err := ioutil.ReadFile(f.Name())
	if err != nil{
		panic(err)
	}
	if len(b) > maxROMSize {
		panic("ROM overread")
	}
	if len(b) < minROMSize {
		panic("rom underread")
	}
	r = b
	return &r
}

func (r *ROM) GetTitle() string {
	return string((*r)[0x0134:0x0142])
}

func (r *ROM) IsGBC() bool {
	return (*r)[0x0143] == 0x80
}

func(r *ROM) LicenseeCode() [2]byte {
	return [2]byte{(*r)[0x0144], (*r)[0x0145]}
}

func(r *ROM) IsSuperGB() bool {
	return (*r)[0x0146] == 0x03
}
