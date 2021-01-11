package mem

type RAM [0xFFFF]byte

func (r *RAM) doWrite(addr uint16, data []byte) {
	//TODO: pass off responsibility to the different memory regions
	//i.e. framebuffer, OAM, etc...
	for i, b := range data {
		offset := uint16(i)
		r[addr+offset] = b
	}
}

//WriteByte writes the given byte to the given address
// the address is expected to be in big-endian, e.g. top of VRAM is uint16 of 9FFF
func (r *RAM) WriteByte(addr uint16, val byte) {
	r.doWrite(addr, []byte{val})
}

// ReadByte is a wrapper over the internal memory
// this allows the various memory controllers to control their regions
// e.g. hardware IO registers, cartridge RAM, OAM, VRAM, etc...
func (r *RAM) ReadByte(addr uint16) byte {
	return r[addr]
}