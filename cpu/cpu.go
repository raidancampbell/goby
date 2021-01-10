package cpu

import (
	"fmt"
	"github.com/raidancampbell/goby/mem"
)

type REG [2]byte

func (r *REG) toUint16() uint16 {
	return (uint16(r[1]) << 8) + uint16(r[0])
}

func (r *REG) fromUint16(u uint16) {
	r[0] = byte(u & 0xFF)
	r[1] = byte(u >> 8)
}

type CPU struct {
	pc                  uint16 // program counter
	sp                  uint16 // stack pointer
	accFlagReg          REG    // AF
	bcREG, deREG, hlREG REG    // BC, DE, HL
	ram                 mem.RAM
	interruptEnabled    bool
}

func GetRAM() *mem.RAM {
	return &c.ram
}

func DryRun() {
	for i := 0; i < 30000; i++ {
		opByte := c.ram.ReadByte(c.pc)
		newOp, ok := table[opByte]
		fmt.Printf("executing opcode %x at location %x, execution number %v\t%s\n", opByte, c.pc, i, newOp.label)
		if !ok {
			panic(fmt.Sprintf("unable to find opcode %x", c.ram.ReadByte(c.pc)))
		}
		newOp.impl()
	}
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

func (c *CPU) popWord() uint16 {
	val := uint16(c.ram.ReadByte(c.sp)) + (uint16(c.ram.ReadByte(c.sp+1)) << 8)
	fmt.Printf("popped %x from stack addr %x\n", val, c.sp)
	c.sp += 2
	return val
}

func (c *CPU) pushWord(word uint16) {
	high := byte(word >> 8)
	low := byte(word & 0xFF)
	c.pushBytes(low, high)
	fmt.Printf("pushed %x to stack addr %x\n", word, c.sp)
}

func (c *CPU) pushBytes(low, high byte) {
	c.sp--
	c.ram.WriteByte(c.sp, high)
	c.sp--
	c.ram.WriteByte(c.sp, low)
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

func InitPCForBootrom() {
	c.pc = 0x0000
}
