package mem

type RAM [0xFFFF]byte
var ram RAM

func (r *RAM) doWrite(addr uint16, data []byte) {
	//TODO: pass off responsibility to the different memory regions
	//i.e. framebuffer, OAM, etc...
	for i, b := range data {
		offset := uint16(i)
		ram[addr+offset] = b
	}
}

func (r *RAM) WriteWord(addr, val uint16) {
	r.doWrite(addr, []byte{byte(val >> 8), byte(val)})
}

func (r *RAM) WriteByte(addr uint16, val byte) {
	r.doWrite(addr, []byte{val})
}

func (r *RAM) ReadByte(addr uint16) byte {
	return r[addr]
}