package cpu

import "github.com/raidancampbell/goby/mem"

type REG [2]byte

func (r REG) toUint16() uint16 {
	return (uint16(r[0]) << 8) + uint16(r[1])
}

type CPU struct {
	pc                  uint16 // program counter
	sp                  uint16 // stack pointer
	accFlagReg          REG    // AF
	bcREG, deREG, hlREG REG    // BC, DE, HL
	ram                 mem.RAM
}

//setFlag sets the given bit in the flag register to the given value (i.e. setFlag can clear a bit)
func (c *CPU) setFlag(flag uint8, isSet bool) {
	if isSet {
		c.accFlagReg[1] |= 1 << flag
	} else {
		var mask byte = ^(1 << flag)
		c.accFlagReg[1] &= mask
	}
}

func (c *CPU) getFlag(flag uint8) bool {
	return (c.accFlagReg[1] & (1 << flag)) > 0
}

const (
	flagZero      = 0x7 //Z
	flagSubtract  = 0x6 //N
	flagHalfCarry = 0x5 //H
	flagCarry     = 0x5 //C
)

var c *CPU

func init() {
	c = &CPU{}
	c.accFlagReg = [2]byte{0x01, 0xb0}
	c.bcREG = [2]byte{0x00, 0x13}
	c.deREG = [2]byte{0x00, 0xd8}
	c.hlREG = [2]byte{0x01, 0x4d}
	c.pc = 0x1000
	c.sp = 0xFFFE
}
